{{range $nsname,$ns := .}}Namespace: {{$nsname}}
  {{range $,$en := $ns.Enumerates}}Enum[{{$en.Name}}] {{if $en.Comment}}//{{$en.Comment}}{{end}}
    {{range $,$e := $en.EnumValue}} {{printf "%-20s%5d" $e.Name $e.Value}} {{if $e.Comment}}//{{$e.Comment}}{{end}}
    {{end}}
{{end}}{{end -}}
