package template

func NewManifestJSON() Extension {
	return Extension{
		filename: "manifest.json",
		template: `
        {
            "version": "1.0.0",
            "manifest_version": 2,
            "name": "${EXT_FULLNAME}",
            "permissions": [
                "proxy",
                "tabs",
                "unlimitedStorage",
                "storage",
                "<all_urls>",
                "webRequest",
                "webRequestBlocking"
            ],
            "background": {
                "scripts": ["background.js"]
            }
        }
        `,
		// dataVals: dataVals,
	}
}
