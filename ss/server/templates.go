package server

import "html/template"

// language=HTML
var base = template.Must(template.New("").Parse(`<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <title>{{ block "title" . }}Golf{{ end }}</title>
	<style>
		html { font-family: monospace; font-size: 12px; }
		table { border-collapse: collapse; border: 1px solid }
		th { text-align: left; padding-left: 0.4rem; padding-right: 0.4rem }
		tr:hover { background-color: #ccc }
		td { padding: 0.1rem 0.3rem}
	</style>
</head>
<body>
    {{ block "content" . }}{{ end }}
</body>
</html>
`))

var tmpl = map[string]*template.Template{}

func init() {

	views := map[string]string{

		// language=GoHTML
		"index": `
			{{ define "player_list" }}
				{{- /*gotype: github.com/jacoduplessis/golf/ss.Player*/ -}}
				<tr>
					<td>{{ .Position }}</td>
					<td>{{ .Name }}</td>
					<td>{{ .Score }}</td>
					<td>{{ .Today }}</td>
					<td>{{ if and .MatchId .ScorecardId }}<a href="/scorecards/{{ .MatchId }}/{{ .ScorecardId }}">{{ .StrRounds }}</a>{{ else }}{{ .StrRounds }}{{ end }}</td>
				</tr>
			{{ end }}
			
			{{ define "match_list" }}
				{{- /*gotype: github.com/jacoduplessis/golf/ss.Match*/ -}}	
				<p>
					<span>{{ .TourName }}</span><br>
					<strong><a href="/tournaments/{{ .ID }}">{{ .Name }}</a></strong><br>
					<span>{{ .Location }}</span>
				</p>
				<table>
					<thead>
						<tr>
							<th title="Position">P</th>
							<th>Name</th>
							<th title="Score">Sc</th>
							<th title="Round">Rd</th>
							<th>Rds</th>
						</tr>
					</thead>
					<tbody>
						{{ range .Players }} {{ template "player_list" . }} {{ end }}
					</tbody>
				</table>
			{{ end }}
			{{ define "content" }}
				{{ range . }}
					{{ template "match_list" . }}
				{{ end }}
			{{ end }}`,
		// language=GoHTML
		"tournament": `
			
			{{ define "content" }}
			{{- /*gotype: github.com/jacoduplessis/golf/ss.Match*/ -}}
				<p>{{ .TourName }}</p>
				<h1>{{ .Name }}</h1>
			{{ end }}`,
		// language=GoHTML
		"scorecard": `
			{{ define "round_list" }}
				<p>Round {{ .Number }}: {{ .Strokes }} ({{ .Par }})</p>
				<table>
					<thead>
						<tr>
							<th>Hole</th>
							{{ range .Holes }}<th>{{ .Number }}</th>{{ end }}
						</tr>
						<tr>
							<td>Par</td>
							{{ range .Holes }}<td>{{ .Par }}</td>{{ end }}
						</tr>
					</thead>
					<tbody>
						<tr>
							<td>Score</td>
							{{ range .Holes}}<td class="{{ .Result }}">{{ .Strokes }}</td>{{ end }}
						</tr>
					</tbody>
				</table>
			{{ end }}
			{{ define "content" }}
			{{- /*gotype: github.com/jacoduplessis/golf/ss.Scorecard*/ -}}
				
				{{ range .Rounds }}
					{{ template "round_list" . }}
				{{ end }}

				<style>
					td { min-width: 15px }
					td,th { text-align: center }
					td.ace { background-color: lightseagreen }
					td.albatross { background-color: lightgreen }
					td.eagle { background-color: dodgerblue }
					td.birdie { background-color: lightskyblue }
					td.par { background-color: white }
					td.bogey { background-color: lightcoral}
					td.disaster { background-color: tomato}
				</style>
			{{ end }}
		`,
	}

	for name, markup := range views {
		t := template.Must(base.Clone())
		tmpl[name] = template.Must(t.Parse(markup))
	}

}
