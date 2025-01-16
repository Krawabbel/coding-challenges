package lcs

type node[T comparable] struct {
	length int
	curr   T
	prev   *node[T]
}

func lcsImpl[T comparable, S ~[]T](s1, s2 S) S {
	m := len(s1)
	n := len(s2)
	c := make([][]*node[T], m+1)
	for i := range c {
		c[i] = make([]*node[T], n+1)
		c[i][0] = new(node[T])
	}

	for j := range c[0] {
		c[0][j] = new(node[T])
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if s1[i-1] == s2[j-1] {
				c[i][j] = &node[T]{c[i-1][j-1].length + 1, s1[i-1], c[i-1][j-1]}
			} else {
				c[i][j] = longer(c[i-1][j], c[i][j-1])
			}
		}
	}

	return backtrace[T, S](c[m][n])
}

func backtrace[T comparable, S ~[]T](n *node[T]) S {
	if n.length == 0 {
		return nil
	}
	return append(backtrace[T, S](n.prev), n.curr)
}

func longer[T comparable](n1, n2 *node[T]) *node[T] {
	if n1.length > n2.length {
		return n1
	} else {
		return n2
	}
}

func lcsstr(s1, s2 string) string {
	return string(lcsImpl([]byte(s1), []byte(s2)))
}
