class FileModel {
    //name of the file
    name: string;
    //Unix timestamp representing UTC time when the file will be removed from the server
    death: number;
    //Compression level (should be plain, gz or xz)
    compression: string;
    constructor(name: string, death: number, compression: string) {
        this.name = name;
        this.death = death;
        this.compression = compression;
    }
}

export default FileModel
