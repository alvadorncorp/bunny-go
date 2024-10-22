# Bunny Go

Bunny Go is a Golang library that provides a simple interface to interact with the [BunnyCDN Storage API](https://docs.bunny.net/docs/storage-api). This project aims to simplify integration with BunnyCDN's storage service for tasks such as file uploads, downloads, deletions, and more.

## Features

- **Simple File Management**: Upload, download, delete, and list files within BunnyCDN storage zones.
- **Optimized for Go**: Written natively in Golang for performance and ease of use.

## Installation

You can install the library using `go get`:

```bash
go get github.com/yourusername/bunny-go
```

## API Documentation

This library currently supports the BunnyCDN Storage API:

- **Upload Files**: Upload local files to BunnyCDN storage. 
- **Download Files**: Retrieve files from BunnyCDN storage. 
- **Delete Files**: Remove files from BunnyCDN storage. 
- **List Files**: List files in a specific directory.

You can refer to the official [BunnyCDN Storage API](https://docs.bunny.net/docs/storage-api) documentation for more details on the available endpoints and usage.

## Contributing

We welcome contributions! If you would like to add more features or fix issues, feel free to open a pull request.

- Fork the repository. 
- Create a new branch (git checkout -b feature-branch). 
- Make your changes and test them. 
- Submit a pull request with a detailed explanation of your changes.

## License

This project is licensed under the Apache License 2.0 - see the LICENSE file for details.
