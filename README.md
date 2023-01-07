# Qkessler-README

Go project to dynamically generate the contents of the [README of my github profile](https://github.com/Qkessler/Qkessler). The way we fill the Github content is by pulling Github API, and then composing the different UI elements by chunking and joining the resulting repos ordered by the last time that there was a push. 

## Developing

If you want to use this or contribute, you can get the Go project by cloning this repo:

```
git clone https://github.com/Qkessler/Qkessler-README
```

and then inspecting the `./readme.sh` file, that contains information on how to set the `GH_ACCESS_TOKEN` environment variable, that contains the access token to Github API. I personally use `pass` to manage my passwords, and if you do too, you can just add an entry called `gh-access-token`.

After that, you should be able to run the `./readme.sh` file, but will be generating the readme dynamically for the Qkessler submodule included in the repo. You can change that and rerun `./readme.sh`.

## Example generated README

You can visit my [main Qkessler page on Github](https://github.com/Qkessler) to see an example of the generated README at all times, but you can also find the content below:

----------------

Hi there! I'm Quique, and I'm currently working as a SDE at Amazon. I'm part of the Notifications team on the Kindle organization, and I work on bringing notifications to customers on different surfaces: iOS, Android and FireOS. On my free time, I like reading, hiking, plants, cats and Open Source Software. Here's my web: [**enriquekesslerm.com**](https://enriquekesslerm.com), where you'll find posts about any of those.

<div align="center">

[Email](mailto:enrique.kesslerm@gmail.com) • [Twitter](https://twitter.com/quique_kessler) • [LinkedIn](https://www.linkedin.com/in/enrique-kessler-martinez/) • [Goodreads](https://www.goodreads.com/user/show/130860665-quique)

</div>

Below you'll find a featured repo, which is **one of my 10th last updated**. Below the repo card, you'll find an ordered list of my repositories and their languages, which shows the languages that I work with the most, and the ones I have been working with as of late, for my personal projects.

<div align="center">
    <a href="https://github.com/Qkessler/consult-project-extra">
        <img src="Qkessler/src/repo-card.svg" alt="Repo card which links to the Repo itself, in Github.">
    </a>
</div>

<div align='center'>

|  **Go**  |  **Rust**  |
| :--: | :--: |
|  [Qkessler/Qkessler-README](https://github.com/Qkessler/Qkessler-README) |   [Qkessler/dyncomp](https://github.com/Qkessler/dyncomp)  |
|  [Qkessler/dyncomp.go](https://github.com/Qkessler/dyncomp.go) |   [Qkessler/santander-ledger](https://github.com/Qkessler/santander-ledger)  |
|  **Emacs Lisp**  |  **C**  |
|  [Qkessler/qk-emacs](https://github.com/Qkessler/qk-emacs) |   [Qkessler/PapsGMP](https://github.com/Qkessler/PapsGMP)  |
|  [Qkessler/consult-project-extra](https://github.com/Qkessler/consult-project-extra) | :small_orange_diamond:  [Qkessler/qmk_firmware](https://github.com/Qkessler/qmk_firmware)  |
| :small_orange_diamond: [Qkessler/emacs-calfw](https://github.com/Qkessler/emacs-calfw) |   |
|  [Qkessler/consult-projectile](https://github.com/Qkessler/consult-projectile) |   |
| :small_orange_diamond: [Qkessler/emacs-wttrin](https://github.com/Qkessler/emacs-wttrin) |   |
| :small_orange_diamond: [Qkessler/emacs-theme-gruvbox](https://github.com/Qkessler/emacs-theme-gruvbox) |   |
| :small_orange_diamond: [Qkessler/emmet-mode](https://github.com/Qkessler/emmet-mode) |   |
|  [Qkessler/dot_files](https://github.com/Qkessler/dot_files) |   |
|  **TeX**  |  **JavaScript**  |
|  [Qkessler/cv](https://github.com/Qkessler/cv) |   [Qkessler/enriquekesslerm.com](https://github.com/Qkessler/enriquekesslerm.com)  |
|  **Java**  |  **Python**  |
|  [Qkessler/AppMusic](https://github.com/Qkessler/AppMusic) |   [Qkessler/CloudQuestions](https://github.com/Qkessler/CloudQuestions)  |
|  [Qkessler/ChatProtocol](https://github.com/Qkessler/ChatProtocol) |   [Qkessler/rotating_background](https://github.com/Qkessler/rotating_background)  |

</div>

