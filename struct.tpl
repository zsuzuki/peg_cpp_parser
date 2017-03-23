{{range $nsname,$ns := .}}Namespace: {{$nsname}}
  {{range $,$st := $ns.StructList}}Struct[{{$st.Name}}] has {{len $st.Variables}} members {{if $st.Comment}}//{{$st.Comment}}{{end}}
    {{range $,$v := $st.Variables}} {{printf "%-20s" $v.Type}}{{if $v.Size}}{{printf "%-20s" (print $v.Name "[" $v.Size "]")}}{{else}}{{printf "%-20s" $v.Name}}{{end}} {{if $v.Comment}}//{{$v.Comment}}{{end}}
    {{end}}
  {{end}}{{end -}}
