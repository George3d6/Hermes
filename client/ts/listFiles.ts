import {displayError, displayMessage} from './views'

const listFiles = () => {
    return new Promise<string>((resolve, reject) => {
        const xhr = new XMLHttpRequest();
        xhr.open("GET", "/get/list/");
        xhr.onload = () => {
            try {
                const res = JSON.parse(xhr.responseText);

                if (res.status === "error") {
                    displayError(res.message);
                } else {
                    displayMessage(res.message);
                }

                reject(Error('Could not load files:' + xhr.statusText));
            } catch (e) {
                resolve(xhr.responseText);
            }
        };
        xhr.send(null);
    });
};

export default listFiles;
