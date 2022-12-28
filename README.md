# goConfigHandler
![Build and Test](https://github.com/KevinCFechtel/utils/actions/workflows/go_build_and_test.yml/badge.svg)

This package provides a "secure" handling of json formatted config files.
It tries to read and save config files encrypted from different locations.
This should make the config handling easier when changeing the deployment of the application from a local server to a cloud environment.

## Methods
### Local File handling
The implementation of handling a local file read the config from a local encrypted file within the filesystem.
If the file ist unencrypted, configHandler read the config file and check if a keyfile already exists.
If a keyfile exists, configHandler decrypt the config file with the existing keyfile and replace the unencryppted config file with the encrypted one.
If a keyfile doesn't exists, configHandler generate a new keyfile and save this together with the encrypted version of the config file in the same location.

## Encryption
To prevent a easy decrypton of the config file if the keyfile and the config file are published, the [encryption] methods use a salt to alter the saved key before decryption or encryption.
This salt is provided by the application and can be a Hostname oder a MAC-Address of the host in which the application is running.
It can also be another value provided by another source, it shouldn't be saved together with the keyfile.

## Whats next?
New methods for handling config files will be added on demand.
