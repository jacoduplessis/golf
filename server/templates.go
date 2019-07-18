package server

import (
	"github.com/jacoduplessis/golf"
	"html/template"
)

func parseTemplate() {
	// language=HTML format=true
	tmpl = template.Must(template.New("leaderboard").Parse(`
	<div style="margin-right: 1rem">	
		<h3>{{ .Tour }} - {{ .Tournament }}</h3>
		<h3>{{ .Course }}{{if .Location}}, {{ .Location }}{{end}}</h3>
		<table>
			<thead>
				<tr>
					<th title="position">Pos</th>
					<th>Par</th>
					<th>Player</th>
					<th title="nationality">Nat</th>
					<th title="current hole">On</th>
					<th title="holes played">Pl</th>
					<th title="this round">Rd</th>
					<th>Rounds</th>
					<th title="number of strokes">Total</th>
				</tr>
			</thead>
			<tbody>
			{{range .Players}}
				<tr>
					<td>{{ .CurrentPosition }}</td>
					<td style="text-align: right">{{ .Total }}</td>
					<td>{{ .Name }}</td>
					<td style="text-align: right">{{ .Country }}</td>
					<td style="text-align: right">{{ .Hole }}</td>
					<td style="text-align: right">{{ .After }}</td>
					<td style="text-align: right">{{ .Today }}</td>
					<td style="padding-left: 1rem">{{range .Rounds}}{{if .}}{{ . }} {{end}}{{end}}</td>
					<td style="text-align: right">{{ .TotalStrokes}}</td>
				</tr>
			{{end}}
			</tbody>
		</table>
	</div>
	`))

	// language=HTML format=true
	tmpl.New("").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
  <title>Golf Leaderboards</title>
  <style>
	html { font-family: monospace }
	table { border-collapse: collapse; border: 1px solid }
	th { text-align: left; padding-left: 0.4rem; padding-right: 0.4rem }
	tr:hover { background-color: #ccc }
	td:nth-of-type(3) { padding-left: 10px }
</style>
</head>
<body>
	<div style="display: flex; flex-flow: row wrap">
	{{range .Leaderboards}}{{template "leaderboard" .}}{{end}}
	</div>
	<p>Click on any country code to highlight all players from that country.</p>
	<p><a href="/?format=json">Get this data as JSON.</a></p>
	<p>View <a href="/news">news</a> or <a href="/results/">past results</a>.</p>
	<script>
		(function(){
			const cells = document.querySelectorAll('td:nth-of-type(4)') // country code
			cells.forEach(el => {
				el.addEventListener('click', event => {					
					cells.forEach(e => e.parentElement.style.backgroundColor = '') // reset all
					cells.forEach(e => {
						if (e.textContent === event.target.textContent) e.parentElement.style.backgroundColor = 'yellow' // highlight all from same country
					})	
				})
			})
		})()
	</script>
</body>
</html>
`)

	newsTemplate = template.Must(template.New("item").
		Funcs(map[string]interface{}{
			"URLize": golf.URLize,
		}). // language=HTML
		Parse(`
			<div class="tweet">
				<p><strong>{{ .UserName }} (@{{ .UserHandle}})</strong> &middot; <time datetime="{{ .ISOTime }}">{{ .RelativeTime }}</time></p>
				
				<p>{{ URLize .Content }}</p>
				{{ if (and .ImageURL (not .Video))}}
				<img src="{{ .ImageURL }}">
				{{end}}

				{{ if .Video }}
				<video controls poster="{{ .VideoThumbnail }}">
					<source src="{{ .VideoSource }}" type="video/mp4">
				</video>
				{{ end }}
			</div>
	`))

	// language=HTML format=true
	newsTemplate.New("").Parse(`
<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
		<title>Golf News</title>
		<style>
			img,video {
				max-width: 100%;
			}

			.tweet {
				margin-top: 1rem;
				border: 2px solid #ccc;
				padding: 1rem;
				background-color: #fff;
			}
		</style>
	</head>
	<body style="max-width: 650px;margin: 0 auto; background-color: #eee">		
		
		<div style="text-align: center; margin-bottom: 3rem">
			<h1>Golf News</h1>
			<p>from Twitter</p>
			<p><a href="/">leaderboards</a></p>
		</div>
		
		
		{{ range . }}{{template "item" . }}{{end}}
		
		<p><a href="/">leaderboards</a></p>
	</body>
</html>
`)

	// language=HTML
	videoTemplate = template.Must(template.New("item").Parse(`
			<div class="video">
				<h3>{{ .Title }}</h3>
				<p>{{ .Description }}</p>

				{{ if .SRC }}
				<video controls poster="{{ .ThumbnailSRC }}">
					<source src="{{ .SRC }}" type="video/mp4">
				</video>
				{{ end }}
			</div>
`))

	// language=HTML
	videoTemplate.New("").Parse(`
<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
		<title>Golf Videos</title>
		<style>
			img,video {
				max-width: 100%;
			}

			.video {
				margin-top: 1rem;
				border: 2px solid #ccc;
				padding: 1rem;
				background-color: #fff;
			}
		</style>
	</head>
	<body style="max-width: 900px;margin: 0 auto; background-color: #eee">		
		
		<div style="text-align: center; margin-bottom: 3rem">
			<h1>Golf Videos</h1>
			<p><a href="/">leaderboards</a></p>
		</div>

		{{ range . }}{{template "item" . }}{{end}}

		<p><a href="/">leaderboards</a></p>
	</body>
</html>
`)

}
