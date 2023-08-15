package git_wrapper

type Branch struct {
	name string
}

func NewBranch(name string) Branch {
	return Branch{
		name: name,
	}
}
