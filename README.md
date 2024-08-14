# Akhanda

Akhanda is a Go program that generates and verifies SHA256 checksums for files in a directory. It supports parallel processing and displays a progress bar while it's working.

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/majorsahab/akhanda.git
    ```

2. Navigate to the project directory:

    ```sh
    cd akhanda
    ```

3. Build the project:

    ```sh
    go build
    ```

## Usage

Akhanda supports two actions: `generate` and `verify`.

### Generate checksums

To generate checksums for all files in a directory, use the `generate` action. By default, Akhanda will process the current directory and save the checksums to a file named `checksums.txt`.

```sh
./akhanda -action=generate -directory=/<path>/<to>/<dir> -checksumFile=<mychecksums_file>
```

### Verify checksums

To verify checksums for all files in a directory, use the `verify` action. Akhanda will read the checksums from the specified checksum file and compare them to the current checksums of the files in the directory.

```sh
./akhanda -action=verify  -checksumFile=<mychecksums_file>
```

## Contributing

Contributions are welcome! Please feel free to submit a pull request.

## License

This project is licensed under the MIT License. See the COPYING file for details.
