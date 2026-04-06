# KSynth

KSynth is a small terminal-first synthesizer and loop sequencer written in Go.

It is built for direct play: open a terminal, press keys, switch timbre, record a loop, inspect active voices, and export the result to a `.wav` file.

## Philosophy

KSynth is intentionally simple.

- It favors immediate interaction over a heavy UI.
- It treats the computer keyboard like a playable instrument.
- It keeps the sound engine understandable enough to modify.
- It is meant for sketching riffs, hooks, backing loops, and ideas quickly.

This is not trying to be a full DAW. It is a compact live-play synth with enough sequencing to build a musical idea fast.

## Features

- Play notes directly from the keyboard
- Three sound modes: `sine`, `piano`, `organ`
- Four sequencer slots for looping note patterns
- Octave shifting while playing
- Sustain toggle on the last triggered voice
- Voice inspection and manual voice removal
- Live audio export to `.wav`

## Run

Run from the project root:

```powershell
go run ./cmd/ksynth
```

Start with a specific sound:

```powershell
go run ./cmd/ksynth -sound organ
go run ./cmd/ksynth -sound piano
go run ./cmd/ksynth -sound sine
```

Show help:

```powershell
go run ./cmd/ksynth -h
```

## Keyboard Mapping

### Notes

The playable note row is:

- `a` = `C4`
- `w` = `C#4`
- `s` = `D4`
- `e` = `D#4`
- `d` = `E4`
- `f` = `F4`
- `t` = `F#4`
- `g` = `G4`
- `y` = `G#4`
- `h` = `A4`
- `u` = `Bb4`

This gives a compact piano-style layout across the keyboard.

### Sound Modes

- `7` = switch to `sine`
- `8` = switch to `piano`
- `9` = switch to `organ`

### Sequencer

- `1` = select slot 1
- `2` = select slot 2
- `3` = select slot 3
- `4` = select slot 4
- `Space` = start or stop recording on the selected slot

### Performance Controls

- `+` = octave up
- `-` = octave down
- `j` = toggle sustain on the last triggered voice

### Voice Tools

- `v` = print active voices
- `k` = remove a voice by numeric voice id

### Export / Exit

- `r` = export live output to `.wav`
- `Esc` = stop all sequences and quit

## How Sequencing Works

KSynth has four sequencer slots.

Basic flow:

1. Press `1`, `2`, `3`, or `4` to choose a slot.
2. Press `Space` to arm recording.
3. Play notes on the keyboard.
4. Press `Space` again to stop recording.
5. That slot immediately starts looping.

Notes:

- Only one slot can be actively recording at a time.
- Starting recording on a slot clears that slot first.
- If you stop recording without playing notes, that slot is cleared.
- Recorded notes are turned into a loop and mixed with live playback.

## Sound Modes

### Sine

Clean and plain. Useful for testing pitch and timing.

### Piano

Brighter attack with a more percussive body. Better for melodic playing.

### Organ

Steadier sustained harmonic tone. Better for long notes and layered loops.

## Voice Inspection

Press `v` to see only active voices.

Example:

```text
Active voices:
[3] C4 key=a freq=261.63Hz sustained=false mode=organ
[5] G4 key=g freq=392.00Hz sustained=true mode=organ
```

This makes `k` usable, because you can identify a sounding voice before removing it.

Press `k`, then enter the numeric voice id shown in the list.

## Export Audio

Press `r` while KSynth is running.

You will be prompted for a duration in seconds. KSynth then writes the current live output to a `.wav` file in the current directory.

Example output name:

```text
recording-20260407-153000.wav
```

The export captures whatever is audible during that time window:

- live notes you play
- active sequencer loops
- current sound mode output

## Build

Build a local Windows binary:

```powershell
go build -o ksynth.exe ./cmd/ksynth
```

Build Windows release archives into `dist/` with GoReleaser:

```powershell
goreleaser release --snapshot --clean
```

Publish a tagged GitHub release:

```powershell
git tag v0.1.0
git push origin v0.1.0
```

Pushing the tag triggers the GitHub Actions release workflow, which runs GoReleaser and publishes the release assets automatically.

## Project Layout

Current structure is intentionally small:

- `cmd/ksynth` contains the terminal app entrypoint
- `internal/audio` contains the synth engine, stream output, waveforms, envelopes, export, and sequencing glue
- `internal/sequencer` contains loop recording and sequence generation

## Limitations

- This is a terminal app, so input is constrained by terminal keyboard handling.
- The synth currently uses one shared engine mode, so changing sound mode affects current playback behavior globally.
- It is designed for quick sketches, not deep editing or arrangement workflows.

## Good First Session

Try this:

1. Run `go run ./cmd/ksynth -sound organ`
2. Press `1`
3. Press `Space`
4. Play a simple rhythm on `a`, `g`, and `f`
5. Press `Space` again to loop it
6. Press `8` to switch to piano
7. Play a lead over the loop
8. Press `r` and export 10 seconds

That is the intended KSynth workflow: fast idea capture, not setup-heavy production.
