package bham

var Filters = []FilterHandler{
	FilterHandler{
		Trigger: ":javascript",
		Open:    `<script type="text/javascript">`,
		Close:   "</script>",
		Handler: Transformer(func(s string) string { return s }),
	},
	FilterHandler{
		Trigger: ":css",
		Open:    `<style>`,
		Close:   "</style>",
		Handler: Transformer(func(s string) string { return s }),
	},
}

type FilterHandler struct {
	Trigger     string
	Open, Close string
	Handler     Transformer
}

type Transformer func(string) string

func shortHandOpen(sh string) token {
	output := token{
		purpose: pse_tag,
		content: "<script type=\"text/javascript\">",
	}
	return output
}

func shortHandClose(sh string) token {
	output := token{
		purpose: pse_tag,
		content: "</script>",
	}
	return output
}
