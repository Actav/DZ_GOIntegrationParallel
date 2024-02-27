package generators

func ParenthesisString(n int) []string {
	var result []string
	generate("", n, n, &result)
	return result
}

func generate(current string, open, close int, result *[]string) {
	if open == 0 && close == 0 {
		*result = append(*result, current)
		return
	}

	if open > 0 {
		generate(current+"(", open-1, close, result)
	}
	if close > open {
		generate(current+")", open, close-1, result)
	}
}
