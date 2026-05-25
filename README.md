# devtag

Play an audio clip of your choice every time you run `git push`. macOS only.

Once set up, every `git push` from any repo, any shell, any IDE will trigger the audio.

## Install

```sh
brew install kelvinxiao01/devtag/devtag
```

If this is your first time tapping the repo, Homebrew will do that automatically. To tap manually:

```sh
brew tap kelvinxiao01/devtag
brew install devtag
```

## Setup

Two commands and you're done.

**1. Pick the audio file to play:**

```sh
devtag set /path/to/your/audio.mp3
```

Supported formats are whatever macOS's `afplay` accepts: `.mp3`, `.wav`, `.aiff`, `.m4a`, `.aac`, `.caf`, `.au`, Apple Lossless. (Not supported: FLAC, OGG, OPUS.)

Try it with a built-in system sound to sanity-check:

```sh
devtag set /System/Library/Sounds/Glass.aiff
devtag play
```

**2. Wire up the global git hook:**

```sh
devtag install
```

That's it. Try `git push` in any repo with a remote — audio should play as the push starts.

## Commands

| Command | What it does |
|---|---|
| `devtag set <path>` | Save the audio file to play on push. |
| `devtag get` | Print the currently configured audio path. |
| `devtag play` | Play the configured audio (also used internally by the git hook). |
| `devtag install` | Install the global git `pre-push` hook. |
| `devtag install --force` | Overwrite an existing `core.hooksPath` setting. |
| `devtag uninstall` | Remove the hook and unset `core.hooksPath`. |
| `devtag --help` | Show usage. |

## How it works

`devtag install` does two things:

1. Writes a `pre-push` script to `~/.devtag/hooks/pre-push`.
2. Runs `git config --global core.hooksPath ~/.devtag/hooks` so git uses that directory for hooks in every repo on the machine.

The hook runs `devtag play` in the background and immediately exits 0, so audio never delays or aborts your push.

## Heads-up: per-repo hooks

While `devtag` is installed, git uses `~/.devtag/hooks/` for hooks instead of each repo's `.git/hooks/`. If you (or tools like Husky, pre-commit, lefthook) rely on per-repo hooks, those will not fire.

To check whether you have any per-repo hooks that would be silenced, scan your workspace:

```sh
find ~/Code ~/projects -maxdepth 6 -type d -name hooks -path "*/.git/hooks" 2>/dev/null \
  | while read d; do
      files=$(ls "$d" 2>/dev/null | grep -v '\.sample$' || true)
      [ -n "$files" ] && { echo "REPO: $(dirname $(dirname $d))"; echo "$files" | sed 's/^/  /'; }
    done
```

Swap `~/Code ~/projects` for wherever you keep your repos. Anything that prints is a real hook that `devtag install` would silence.

If `devtag install` detects you already have `core.hooksPath` set to something else, it refuses to overwrite unless you pass `--force`.

## Uninstall

**Order matters** — run `devtag uninstall` *before* `brew uninstall`, otherwise the hook and `core.hooksPath` setting get left behind.

```sh
devtag uninstall
brew uninstall devtag
```

Optionally:

```sh
brew untap kelvinxiao01/devtag    # stop pulling formula updates
rm -rf ~/.devtag                  # remove the saved audio-path config
```

If you uninstalled in the wrong order, clean up by hand:

```sh
git config --global --unset core.hooksPath
rm -rf ~/.devtag
```

## License

MIT — see [LICENSE](LICENSE).
