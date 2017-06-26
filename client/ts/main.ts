import uploadFile from './uploadFile';
import listFiles from './listFiles';
import router from './router';
import File from './file'
import enableAuthenticationForm from './authenticate'

uploadFile('upload_form');
enableAuthenticationForm();

const initFileListView = () => {
    listFiles().then((answer) => {
        document.getElementById('file_list').innerHTML = '';
        const files: Array<File> = [];
        answer.split('#|#').forEach((stringifiedFile) => {
            if(stringifiedFile==='') {
                return;
            }
            const fileElements = stringifiedFile.split('|#|');
            files.push(
                new File(fileElements[0], fileElements[1], parseInt(fileElements[2]))
            );
            return files;
        });
        document.getElementById('file_list').insertAdjacentHTML('beforeend',
        `<p id="close_file_view_holder"><a href="/#!/empty/" id="close_file_view"><i class="close icon big"></i></a></p>`);
        files.forEach((f) => {
            let fileDies: Date = new Date(f.death * 1000);
            let fileDiesString = fileDies.getUTCFullYear() + "-" + fileDies.getUTCMonth() + "-" + fileDies.getUTCDay() + " " + fileDies.getUTCHours() + ":" + fileDies.getUTCMinutes() + ":" + fileDies.getUTCSeconds();
            document.getElementById('file_list').insertAdjacentHTML('beforeend',
            `
            <div class="item file_in_list">
             <i class="huge cloud download middle aligned icon cursor_hover files_icon" id="download_${f.name}" onclick="window.location='/get/file/?file=${f.name}';"></i>
             <i class="huge trash middle aligned icon cursor_hover trashed_icon" id="trash_${f.name}"></i>
              <div class="content" id="${f.name}" class="inline_content">
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
            `
        );
        document.getElementById(`trash_${f.name}`).addEventListener('click', (e) => {
            const xhr = new XMLHttpRequest();xhr.open('GET', '/delete/file/?file=' + f.name);
            xhr.send(null);
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
