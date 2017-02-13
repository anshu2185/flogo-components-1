package splitpath

var jsonMetadata = `{
  "name": "splitpath",
  "version": "0.0.1",
  "description": "Splits a path into separate parts",
  "title": "Split Path",
  "author": "Jan van der Lugt <jvanderl@tibco.com>",
  "homepage": "https://github.com/jvanderl/flogo-components/activity/splitpath",
  "inputs":[
    {
      "name": "input",
      "type": "string"
    },
    {
      "name": "delimiter",
      "type": "string"
    },
    {
      "name": "fixedpath",
      "type": "string"
    }
  ],
  "outputs": [
    {
      "name": "result",
      "type": "string"
    },
    {
      "name": "fixedpath",
      "type": "string"
    },
    {
      "name": "part1",
      "type": "string"
    },
    {
      "name": "part2",
      "type": "string"
    },
    {
      "name": "part3",
      "type": "string"
    },
    {
      "name": "part4",
      "type": "string"
    },
    {
      "name": "part5",
      "type": "string"
    },
    {
      "name": "part6",
      "type": "string"
    },
    {
      "name": "part7",
      "type": "string"
    },
    {
      "name": "part8",
      "type": "string"
    },
    {
      "name": "remainder",
      "type": "string"
    }
  ]
}`