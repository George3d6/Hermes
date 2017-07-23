class FileModel {
    //name of the file
    name: string;
    //Unix timestamp representing UTC time when the file will be removed from the server
    death: number;
    //Compression level (should be plain, gz or xz)
    compression: string;
    //Is the file public ?
    isPublic: boolean;
    //Size of file in KB
    size: number;
    constructor(name: string, compression: string, death: number, isPublic: boolean, size: number = 0) {
        this.name = name;
        this.death = death;
        this.compression = compression;
        this.isPublic = isPublic;
        this.size = size;
    }
}

export default FileModel
