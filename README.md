# ylc (Yeelight CLI)

`ylc` is a command-line interface (CLI) tool designed to control your Yeelight
smart bulbs. With `ylc`, you can manage, control, and configure your  bulbs
effortlessly from your terminal.

## Features

- **Discover new bulbs**: Automatically discover and update known bulbs on your network
- **Control brightness**: Set the brightness level of your bulbs
- **Toggle power**: Turn your bulbs on or off
- **Set RGB color**: Change the color of your bulbs using RGB values
- **Adjust color temperature**: Modify the color temperature of your bulbs
- **Manage bulbs**: List and delete known bulbs

## Installation

You can install `ylc` using one of the following methods.

### Download from GitHub Releases

1. Go to the [releases page](https://github.com/pugkong/ylc/releases)
2. Download the appropriate release for your operating system
3. Extract the downloaded file and move the executable to a directory included
in your system's PATH

### Install using `go install`

If you have Go installed, you can easily install `ylc` with the following command:

```sh
go install github.com/pugkong/ylc@latest
```

## Usage

Below are some common commands you can use with `ylc`.

### Discover Bulbs

Discover new bulbs or update known bulbs on your network:

```sh
ylc discover
```

### List Bulbs

List all known bulbs:

```sh
ylc list
```

### Show Bulb Info

Display detailed information about a specific bulb:

```sh
ylc info [BULB NAME]
```

### Set Brightness

Set the brightness level of a bulb:

```sh
ylc bright [BULB NAME] [BRIGHTNESS]
```

- `[BRIGHTNESS]` should be a value between 1 and 100.

### Toggle Power

Toggle the power state of a bulb:

```sh
ylc power [BULB NAME]
```

### Set RGB Color

Set the RGB color of a bulb:

```sh
ylc rgb [BULB NAME] [COLOR]
```

- `[COLOR]` should be a hexadecimal value (e.g., `ff0000` for red).

### Set Color Temperature

Set the color temperature of a bulb:

```sh
ylc temperature [BULB NAME] [TEMPERATURE]
```

- `[TEMPERATURE]` should be a value between 1700 and 6500.

### Delete Bulb

Delete a bulb from the known bulbs list:

```sh
ylc delete [BULB NAME]
```

## Command Options

Many commands in `ylc` support additional options:

- `--bg`: Apply the command to the background light of the bulb
- `--effect`, `-e`: Set the effect for the command (`smooth` or `sudden`)
- `--duration`, `-d`: Set the duration of the effect in milliseconds

For example, to set the brightness of a bulb with a smooth effect over 1000 milliseconds:

```sh
ylc bright [BULB NAME] [BRIGHTNESS] --effect smooth --duration 1000
```

## Development

### Building from Source

1. Clone the repository:

    ```sh
    git clone https://github.com/pugkong/ylc.git
    cd ylc
    ```

2. Build the project:

    ```sh
    go build
    ```

3. Move the compiled binary to a directory included in your system's PATH.

### Contributing

Contributions are welcome! Feel free to open issues or submit pull requests
on the [GitHub repository](https://github.com/pugkong/ylc).

## License

`ylc` is licensed under the UNLICENSED License.
See the [LICENSE](https://github.com/pugkong/ylc/blob/master/LICENSE) file for
more information.
