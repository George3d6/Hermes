import uploadFile from './uploadFile';
import listFiles from './listFiles';
import router from './router';
import File from './file'
import enableAuthenticationForm from './authenticate'

uploadFile('upload_form');
enableAuthenticationForm();

router.on('/files/', () => {
    listFiles().then((answer) => {
        const files: Array<File> = [];
        answer.split('#|#').forEach((stringifiedFile) => {
            if(stringifiedFile==='') {
                return;
            }
            const fileElements = stringifiedFile.split('|#|');
            files.push(
                new File(fileElements[0], parseInt(fileElements[1]), fileElements[2])
            );
            return files;
        });
        document.getElementById('file_list').insertAdjacentHTML('beforeend',
        `<p id="close_file_view_holder"><a href="/#!" id="close_file_view"><i class="close icon big"></i></a></p>`);
        files.forEach((f) => {
            let fileDies: Date = new Date(f.death * 1000);
            let fileDiesString = fileDies.getUTCFullYear() + "-" + fileDies.getUTCMonth() + "-" + fileDies.getUTCDay() + " " + fileDies.getUTCHours() + ":" + fileDies.getUTCMinutes() + ":" + fileDies.getUTCSeconds();
            document.getElementById('file_list').insertAdjacentHTML('beforeend',
            `
            <div class="item file_in_list">
             <i class="huge file middle aligned icon cursor_hover" onclick="window.location='/get/file/?file=${f.name}';"></i>
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
                <i class="huge trash middle aligned icon cursor_hover" onclick="const xhr = new XMLHttpRequest();xhr.open('GET', '/delete/file/?file=${f.name}');xhr.send(null);"></i>
                </div>
            </div>
            `
        );
        });
    })
        .catch((err) => {
            console.log("Something went horribly wrong: " + err);
        })
}).resolve();

router.on('/', () => {
    document.getElementById('file_list').innerHTML = '';
    document.getElementById('permission_view').style.display = 'none';
}).resolve();
