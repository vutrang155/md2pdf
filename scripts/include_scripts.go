package main 

import (
	"fmt"
	"os"
)

func main() {
	out, _ := os.Create("scripts.go")
	out.Write([]byte("package main\n\n"))
	out.Write([]byte("const MATHJAX_SCRIPT = `\n"))
    out.WriteString(fmt.Sprintf(`
        <script type="text/x-mathjax-config">
        MathJax.Hub.Config({
            TeX: {extensions: ["mhchem.js"]},
            tex2jax: {
            inlineMath: [['$','$'], ['\\(','\\)']],
            displayMath: [ ['$$','$$'], ["\\[","\\]"] ],
            processEscapes: true
            }
        });
        </script>
        <script type="text/javascript"
            src="https://cdnjs.cloudflare.com/ajax/libs/mathjax/2.7.1/MathJax.js?config=TeX-MML-AM_CHTML">
        </script>`))
	out.Write([]byte("`\n"))
}
