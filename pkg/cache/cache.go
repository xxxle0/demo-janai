package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"operarius/internal/models"
	"operarius/pkg/git_wrapper"
	"os"
	"sync"

	"github.com/go-extras/tahwil"
	"github.com/jinzhu/copier"
)

var (
	lruCacheOnce sync.Once
	lruCache     *LRUCache
	lock         sync.Mutex
)

type LRUCache struct {
	Head     *Node            `json:"head"`
	Tail     *Node            `json:"tail"`
	Mapcache map[string]*Node `json:"mapCache"`
	// KB
	CurrentSize uint64 `json:"currentSize"`
	// KB
	SizeLimit uint64 `json:"sizeLimit"`
	// path cache_metadata.json
	Path string `json:"path"`
}

type Node struct {
	Next *Node                       `json:"next"`
	Prev *Node                       `json:"prev"`
	Val  *models.RepositoryConcreate `json:"val"`
}

func LRUCacheSingleton(sizeLimit uint64, rootFolderPath string, path string) *LRUCache {
	lruCacheOnce.Do(func() {
		// check folder /tmp/root exists or not, if not create the new one
		if _, err := os.Stat(rootFolderPath); errors.Is(err, os.ErrNotExist) {
			err := os.MkdirAll(rootFolderPath, os.ModePerm)
			if err != nil {
				log.Println(err)
			}
		}
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			lruCache = RebuildCacheFromDirectory(rootFolderPath, path, sizeLimit)
		} else {
			lruCache = RebuildCacheFromFile(path, sizeLimit)
			if lruCache == nil {
				// remove old cache_metadata file and ignore error
				os.Remove(path)
				log.Println("Reconstruct from file fail fallback to reconstruct from directory")
				lruCache = RebuildCacheFromDirectory(rootFolderPath, path, sizeLimit)
			}
		}
	})
	return lruCache
}

func NewLRUCache(memLimit uint64) *LRUCache {
	return &LRUCache{
		SizeLimit: memLimit,
		Mapcache:  make(map[string]*Node),
	}
}

func (c *LRUCache) Get(key string) *models.RepositoryConcreate {
	lock.Lock()
	defer lock.Unlock()
	if node, found := c.Mapcache[key]; found {
		repository := node.Val
		c.removeNode(key)
		c.addHead(key, repository)
		c.Sync(c.Path)
		return repository
	}
	return nil
}

func (c *LRUCache) Put(key string, val *models.RepositoryConcreate) {
	lock.Lock()
	defer lock.Unlock()
	if node, found := c.Mapcache[key]; found {
		repository := node.Val
		c.removeNode(key)
		c.addHead(key, repository)
	} else {
		c.addHead(key, val)
	}
	// for c.CurrentSize > c.SizeLimit {
	// 	log.Println("Remove worktree")
	// 	if c.Tail == nil {
	// 		break
	// 	}
	// 	c.Tail.Val.RootRepository.RemoveRepository()
	// 	log.Println(c.Tail.Val)
	// 	c.removeTail()
	// }
	c.Sync(c.Path)
}

func (c *LRUCache) removeNode(key string) {
	node, found := c.Mapcache[key]
	if !found || node == nil {
		return
	}
	nNext := node.Next
	nPrev := node.Prev
	if nNext != nil {
		nNext.Prev = nPrev
	}
	if nPrev != nil {
		nPrev.Next = nNext
	}
	if c.Head == node {
		c.Head = c.Head.Next
	}
	if c.Tail == node {
		c.Tail.Prev = nil
	}
	// Remove cache from in-memory
	delete(c.Mapcache, key)
	c.CurrentSize -= node.Val.Size
}

func (c *LRUCache) RemoveNodeLock(key string) {
	lock.Lock()
	defer lock.Unlock()
	node, found := c.Mapcache[key]
	if !found || node == nil {
		return
	}
	// Remove root source from node disk
	node.Val.RootRepository.RemoveRepository()
	c.removeNode(key)
	c.Sync(c.Path)
}

func (c *LRUCache) addHead(key string, val *models.RepositoryConcreate) {
	if c.Head != nil {
		newNode := &Node{
			Next: c.Head,
			Val:  val,
		}
		c.Head.Prev = newNode
		c.Head = newNode
	} else {
		newNode := &Node{
			Val: val,
		}
		c.Head = newNode
		c.Tail = newNode
	}
	c.Mapcache[key] = c.Head
	c.CurrentSize += val.Size
}

func (c *LRUCache) removeTail() {
	if c.Tail != nil {
		log.Println(c.Tail)
		c.CurrentSize -= c.Tail.Val.Size
		if c.Tail.Prev != nil {
			c.Tail = c.Tail.Prev
			c.Tail.Next = nil
		}
	}
}

func (c *LRUCache) Sync(path string) {
	log.Println("==== Sync Cache ====", path)
	v, err := tahwil.ToValue(c.Head)
	if err != nil {
		log.Println(err)
	}
	json, err := json.Marshal(v)
	if err != nil {
		log.Println("unmarshalling fail", err)
	}
	err = os.WriteFile(path, json, 0644)
	if err != nil {
		log.Println("write file error", err)
	}
}

func RebuildCacheFromFile(path string, sizeLimit uint64) *LRUCache {
	log.Println("Reconstruct cache from file")
	lruCache := &LRUCache{
		Mapcache:  make(map[string]*Node),
		SizeLimit: sizeLimit,
		Path:      path,
	}
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Println("read file error", err)
		return nil
	}
	data := &tahwil.Value{}
	err = json.Unmarshal(bytes, data)
	if err != nil {
		log.Println("unmarshalling fail", err)
		return nil
	}
	head := &Node{}
	err = tahwil.FromValue(data, head)
	if err != nil {
		log.Println("unmarshalling fail", err)
		return nil
	}
	cur := head
	for cur != nil {
		lruCache.Mapcache[fmt.Sprintf("%s-%s", cur.Val.Provider, cur.Val.ProviderInternalId)] = cur
		lruCache.CurrentSize += cur.Val.Size
		if cur.Next == nil {
			break
		}
		cur = cur.Next
	}
	lruCache.Head = head
	lruCache.Tail = cur
	return lruCache
}

func RebuildCacheFromDirectory(path string, fileMetaPath string, sizeLimit uint64) *LRUCache {
	lruCache = &LRUCache{
		SizeLimit: sizeLimit,
		Mapcache:  make(map[string]*Node),
		Path:      fileMetaPath,
	}
	log.Printf("Reconstruct the cache from directory %s", path)
	// walk through directory and rebuild the cache
	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		return nil
	}
	for _, dir := range dirs {
		dirName := dir.Name()
		dirPath := path + "/" + dirName
		rootRepository, err := git_wrapper.Load(dirPath)
		if err != nil {
			continue
		}
		repository := &models.Repository{
			Name:           dirName,
			RootRepository: rootRepository,
			Dest:           dirPath,
			IsRoot:         true,
		}

		var concrete models.RepositoryConcreate
		if copier.Copy(&concrete, repository); err == nil {
			lruCache.Put(dirName, &concrete)
		}
	}
	return lruCache
}
