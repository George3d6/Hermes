import uploadFile from './uploadFile';
import listFiles from './listFiles';
import router from './router';
import File from './file'
import enableAuthenticationForm from './authenticate'
import {displayError, displayMessage} from './views'

uploadFile('upload_form');
enableAuthenticationForm();

const initFileListView = () => {
    listFiles().then((answer) => {
        document.getElementById('file_list').innerHTML = '';
        const files: Array<File> = [];
        answer.split('#|#').forEach((stringifiedFile) => {
            if (stringifiedFile === '') {
                return;
            }
            const fileElements = stringifiedFile.split('|#|');
            files.push(
                new File(fileElements[0], fileElements[1], parseInt(fileElements[2]), <any>fileElements[3])
            );
            return files;
        });
        document.getElementById('file_list').insertAdjacentHTML('beforeend',
            `<p id="close_file_view_holder"><a href="/#!/empty/" id="close_file_view"><i class="close icon big"></i></a></p>`);
        files.forEach((f) => {
            const port_str : string = window.location.port === undefined ? '' :  ":" + window.location.port;
            const hostName: string = window.location.protocol + "//" + window.location.hostname + port_str;
            let fileDies: Date = new Date(f.death * 1000);
            let fileDiesString = fileDies.getUTCFullYear() + "-" + fileDies.getUTCMonth() + "-" + fileDies.getUTCDay() + " " + fileDies.getUTCHours() + ":" + fileDies.getUTCMinutes() + ":" + fileDies.getUTCSeconds();
            const urlFileName: string = encodeURIComponent(f.name);
            const fileLink: string = `${hostName}/get/file/?file=${urlFileName}`;
            document.getElementById('file_list').insertAdjacentHTML('beforeend',
            `
            <div class="item file_in_list">
             <input id="link_holder_${urlFileName}" style="display: none;" tabindex="1" autocomplete="off" style="width:1px !important; height:1px !important" type="text"></input>
             <div class="file_control_icon_holder">
             <i class="huge cloud download middle aligned icon list_icon cursor_hover files_icon" id="download_${urlFileName}" onclick="window.location='/get/file/?file=${f.name}';" title="Download file"></i>
             <i class="huge trash middle aligned icon list_icon cursor_hover trashed_icon" id="trash_${urlFileName}" title="Delete file"></i>
             <i class="huge copy middle aligned icon list_icon cursor_hover copy_icon" id="copy_${urlFileName}" title="Copy link to clipboard"></i>
             </div>
             <div class="content" id="${urlFileName}">
                <p>
                    ${f.name}
                </p>
                <p>
                    Valid until: ${fileDiesString}
                </p>
                <div>
                    Compression: ${f.compression}
                </div>
                </div>
            </div>
            <br>
            `
            );
            const copyLinkElement = (<any>document.getElementById(`link_holder_${urlFileName}`));
            copyLinkElement.value = fileLink;
            document.getElementById(`copy_${urlFileName}`).addEventListener('click', (e) => {
                copyLinkElement.style.display = "inline";
                copyLinkElement.select();
                console.log((<any>document.getElementById(`link_holder_${urlFileName}`)).value);
                document.execCommand('copy');
                //clear selection
                if ((<any>document).selection) {
                    (<any>document).selection.empty();
                } else if (window.getSelection) {
                    window.getSelection().removeAllRanges();
                }
                copyLinkElement.style.display = "none";
            });

            document.getElementById(`trash_${urlFileName}`).addEventListener('click', (e) => {
                const xhr = new XMLHttpRequest(); xhr.open('GET', '/delete/file/?file=' + urlFileName);
                xhr.send(null);
                xhr.onload = () => {
                    console.log(xhr.responseText);
                    const res = JSON.parse(xhr.responseText);
                    if (res.status === "error") {
                        displayError(res.message);
                    } else {
                        displayMessage(res.message);
                    }
                };
                console.log("Working");
                initFileListView();
            });
        });
    })
        .catch((err) => {
            console.log("Something went horribly wrong: " + err);
        })
}

window['initFileListView'] = initFileListView;

router.on('/files/', () => {
    initFileListView();
}).resolve();

router.on('/', () => {
    initFileListView();
}).resolve();

router.on('/empty/', () => {
    document.getElementById('file_list').innerHTML = '';
    console.log("HERE:", document.getElementById('file_list').innerHTML);
}).resolve();
