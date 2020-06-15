# Generic Go Template funcs

Usage: 

```go
import ( 
    "github.com/gpmd/gotemplate"
)

func main() {
    gotemplate.Template("{{(xml_decode .).mycode}}", "<?xml version=\"1.0\"?><mycode>MYCODE</mycode>")
}
```
