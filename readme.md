# Save FTP folder to zip (Previously Nitroserv Palworld FTP save to .zip)

This project allow to copy a folder from an FTP server to a local .zip file.

It exist because of Palworld catastrophic characters saves corruption bug and requirement to contact Nitroserv for server save retreiving.

Note : This project also serve as training, it's not perfect, can contain bug, use at your own risks.

## Dependencies (For compilation)

- Go 1.26.6
- github.com/jlaffaye/ftp : FTP Package from jlaffaye

## Steps

### 1 - Note required informations

Basicaly , it's `.env.sample` file content.

- FTP server Address
- FTP server Port
- FTP server Username
- FTP server Password
- FTP server folder to save path
- Local absolute folder path where .zip will be created

### 2 - Compile project (if not compiled)

Use the Makefile with `MAKE`.

Compiled project will appear in the current folder under `/dist/[win/linux]/`

### 3 - Set environment variables

In `.env.sample` , you will have all the required environment variables.

Set them before launching the program.

### 4 - Start

The save will be created in given folder.