# Obtainium Helper

This is a helper template for [Obtainium](https://github.com/ImranR98/Obtainium). To sum this project up, it is used to get the APK(s) that is(are) not possible for the built-in functionality of Obtainium to get.

## How to setup the settings.json file

-   `to_download` - This is the list of APKs that you want to download. You can add as many as you want.
    -   `APK` - This is the name of the APK that you want to download. Its not required that you have the same name as the APK.
        -   `name` - This is the name of the APK file. It will be auto added when the APK is downloaded.
        -   `type` - This is how the APK will be located. It can be either `enumeration`, `web`, or `artifact` (github artifact).
        -   `identifier` - This is a list of things to identify the APK with.
            -   `latest` - This is the latest identifier of the APK.
            -   `enum_limit` - This is the enumeration limit. This is only used for `enumeration`.
            -   `increment_limit` - This is the increment limit. This is only used for `enumeration`.
            -   `pattern` - This is the [regex pattern](https://en.wikipedia.org/wiki/Regular_expression) for grabbing the identifier. Keep in kind this parses the APK URL.
        -   `url` - This is the URL to get the APK.
        -   `agent` - This is the user agent to use when downloading the APK. (Optional)
        -   `patterns` - This is a list of patterns to use to get the APK.
        -   `url_encoded` - This is the encoding of the URL. (Optional)
-   `wait_time` - This is the time in seconds that the script will wait before re-checking for updates.

### More details:

More info on some confusing parts of the settings.json file.

#### Enumeration:

`enumeration` - so far only increments the version.\
`enum_limit` - how many times it will be incremented. For example, if the latest version is 1.0 and the enum_limit is 500, it will increment to 1.500.\
`increment_limit` - how much it will increment by. For example, if the latest version is 1.0 and the increment_limit is 999, it will increment to 1.999.

#### Web:

`patterns` - a list of patterns to use to get the APK. So if you set `url` to a direct APK link you don't need to add patterns. But if there is a website that contains a button or a link somewhere you can add a regex pattern to get the link. The reason why its an array is because if you need to get a link inside the link that the pattern before-hand got, you can add another pattern to get that link.

#### Artifact:

`url` - you should go to the workflow that its running and filter them by what branch you want, then use that link. Link example: `https://github.com/rebelonion/Dantotsu/actions/workflows/beta.yml?query=branch%3Adev`\
`pattern` - you should just use this pattern for now: `\\/([a-zA-Z0-9]+)\\.zip` because there is now way to get the file version without downloading the whole thing, which would use API calls and would be slow.\
`patterns` - the first pattern should be the name of the artifact thats in the action, the second pattern should be the name of the APK thats in the zip file.

#### Other:

`url_encoded` is if the URL is encoded. For example, if the URL a bunch of weird symbols that start with `%` then you need to set this to `true`.\
`agent` is the user agent to use when downloading the APK. This is optional and only needed if the website requires a user agent to download the APK.

### Example:

```json
{
	"to_download": {
		"Codex": {
			"name": "Codex_2.623.apk",
			"type": "enumeration",
			"identifier": { "latest": "2.623", "enum_limit": 500, "increment_limit": 999, "pattern": "\\d+(?:\\.\\d+)+" },
			"url": "https://cdn.codex.lol/public/Codex_",
			"agent": "Codex Android",
			"patterns": null,
			"url_encoded": false
		},
		"Dantotsu": {
			"name": "app-google-alpha.apk",
			"type": "artifact",
			"identifier": {
				"latest": "f0f2d2aba6b71f9649c29cb455b4953be0459fa879643ab4f988b56bc290c564",
				"enum_limit": 0,
				"increment_limit": 0,
				"pattern": "\\/([a-zA-Z0-9]+)\\.zip"
			},
			"url": "https://github.com/rebelonion/Dantotsu/actions/workflows/beta.yml?query=branch%3Adev",
			"agent": null,
			"patterns": ["Dantotsu", "app-google-alpha\\.apk"],
			"url_encoded": false
		},
		"TikTok": {
			"name": "34.9.5-universal.apk",
			"type": "web",
			"identifier": { "latest": "34.9.5", "enum_limit": 0, "increment_limit": 0, "pattern": "\\/(\\d+(?:\\.\\d+)+)" },
			"url": "https://store.kde.org/p/1515346/loadFiles",
			"agent": null,
			"patterns": ["https://files04\\.pling\\.com/api/files/download/j/[^/]+/\\d+(?:\\.\\d+)+-universal\\.apk"],
			"url_encoded": true
		},
		"TikTok Plugin": {
			"name": "1.41_plugin.apk",
			"type": "web",
			"identifier": { "latest": "1.41", "enum_limit": 0, "increment_limit": 0, "pattern": "\\d+(?:\\.\\d+)+" },
			"url": "https://tiktokmodcloud.com/ttplugin.html",
			"agent": null,
			"patterns": ["https://oxy\\.cloud\\/[^\\s\"']+", "https://loader\\.oxy\\.st\\/get\\/[^\\s\"']+"],
			"url_encoded": false
		}
	},
	"wait_time": 15
}
```
