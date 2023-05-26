# 0xyWD
Proxied webdriver session util package written in GoLang

## Supported capabilities
- Google Chrome
- ~~Firefox~~

## Include in project
`import oxywd "github.com/sik0-o/0xyWD"`

## Usage
```
var proxyURL string

...

chromeCaps := chrome.Capabilities{}

if err := oxywd.SetProxy(&chromeCaps, proxyURL); err != nil {
  return nil, err
}

caps := selenium.Capabilities{}
caps.AddChrome(chromeCaps)

// init selenium session with capabilities
```

## Basic API
### SetProxy(CapsExtAdder, string) error 
Build and add proxy extension with provded proxyURL string in browser capabilities that have `CapsExtAdder interface`.
