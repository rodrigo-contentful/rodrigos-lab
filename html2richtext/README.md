# lightweight html to rich text api

Using an endpoint, it will transform html to richtext via markdown.

html -> markdown -> richtext

# Usage:

```
$ npm install

$ npm start

$ curl --location --request POST 'http://localhost:3000/convert' \
--header 'Content-Type: application/json' \
--data-raw '{
	"html": "<h2>What is Lorem Ipsum?</h2><p><strong>Lorem Ipsum</strong> is <a href='www.test.com'>test link</a>simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.</p></div><div>"
}'
```

# References

Big thanks to Morten Bujordet and team.