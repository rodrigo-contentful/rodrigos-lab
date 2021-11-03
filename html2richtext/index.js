const express = require('express')
const bodyParser = require('body-parser')
const TurndownService = require('turndown')
const { richTextFromMarkdown } = require('@contentful/rich-text-from-markdown')
const asyncHandler = require('express-async-handler')
const app = express()

const turndownService = new TurndownService()

const htmlToMarkdown = (html) => {
    return turndownService.turndown(html)
}

const markdownToRichtext = (markdown) => {
    return richTextFromMarkdown(markdown)
}

const htmlToRichText = (html) => {
    return markdownToRichtext(htmlToMarkdown(html))
}

app.use(bodyParser.json())

app.post('/convert', asyncHandler(async (req, res, next) => {

    const html = req.body.html
    const richtext = await htmlToRichText(html)
    res.send(richtext)    
                   
}))

app.listen(3000, '0.0.0.0')