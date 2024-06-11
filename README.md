# Obtainium Helper

Obtainium Helper is a template for [Obtainium](https://github.com/ImranR98/Obtainium), designed to assist in downloading APKs that are beyond the reach of Obtainium's built-in functionality.

## How to Use

1. **Fork the Repository:**
   Follow the standard procedure to fork this repository to your own GitHub account.
2. **Star the Project:**
   If you find this project helpful, please consider leaving a star!

## Setting Up the `settings.json` File

### Main Configuration

-   **`to_download`**: A list of APKs you want to download. You can add as many as needed.

    -   **`APK`**: The name of the APK to download. It doesn't have to match the actual APK name.
        -   **`name`**: The name for the downloaded APK file.
        -   **`type`**: The method to locate the APK. Options are `enumeration`, `web`, or `artifact` (GitHub artifact).
        -   **`identifier`**: Information to identify the APK.
            -   **`latest`**: The latest identifier of the APK.
            -   **`enum_limit`**: Used only for `enumeration`. Sets the limit for enumeration.
            -   **`increment_limit`**: Used only for `enumeration`. Sets the increment limit.
            -   **`pattern`**: A [regex pattern](https://en.wikipedia.org/wiki/Regular_expression) for parsing the APK URL.
        -   **`url`**: The URL to download the APK from.
        -   **`agent`**: (Optional) The user agent to use when downloading the APK.
        -   **`patterns`**: A list of patterns to locate the APK.
        -   **`url_encoded`**: (Optional) Set to `true` if the URL is encoded.

-   **`wait_time`**: The interval (in seconds) the script will wait before re-checking for updates.

### Detailed Explanations

#### Enumeration

-   **`enumeration`**: Increments the version.
-   **`enum_limit`**: Maximum number of increments. E.g., if the latest version is 1.0 and `enum_limit` is 500, it will increment to 1.500.
-   **`increment_limit`**: Increment step. E.g., if the latest version is 1.0 and `increment_limit` is 999, it will increment to 1.999.

#### Web

-   **`patterns`**: List of patterns to locate the APK. If `url` is a direct APK link, patterns aren't needed. If a website contains a button/link, add regex patterns to extract the link. Multiple patterns can be used to extract nested links.

#### Artifact

-   **`url`**: Use the workflow URL filtered by branch. Example: `https://github.com/rebelonion/Dantotsu/actions/workflows/beta.yml?query=branch%3Adev`.
-   **`pattern`**: Use `\\/([a-zA-Z0-9]+)\\.zip` to extract file names.
-   **`patterns`**: First pattern is the artifact name in the action, the second is the APK name in the zip file.

#### Other

-   **`url_encoded`**: Set to `true` if the URL contains encoded characters (e.g., symbols starting with `%`).
-   **`agent`**: Optional user agent for downloading the APK.

### Example `settings.json`

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
