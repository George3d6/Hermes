
const listFiles = () => {
    return new Promise<string>((resolve, reject) => {
        const xhr = new XMLHttpRequest();
        xhr.open("GET", "/list/");
        xhr.onload = () => {
            if (xhr.status === 200) {
                resolve(xhr.responseText);
            } else {
                reject(Error('Could not load files:' + xhr.statusText));
            }
        };
        xhr.send(null);
    });
};

export default listFiles;
