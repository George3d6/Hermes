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
            const fileElements = stringifiedFile.split('|#|');
            files.push(
                new File(fileElements[0], parseInt(fileElements[1]), fileElements[2])
            );
            return files;
        });
        document.getElementById('file_list').insertAdjacentHTML('beforeend',
        `<p id="close_file_view_holder"><a href="/#!" id="close_file_view"><i class="close icon big"></i></a></p>`);
        files.forEach((f) => {
            let fileDies = f.death;
            document.getElementById('file_list').insertAdjacentHTML('beforeend',
                `
            <div class="item file_in_list">
             <i class="huge file middle aligned icon cursor_hover" onclick="window.location='/get/file/?file=${f.name}';"></i>
              <div class="content" id="${f.name}" class="inline_content">
                <p>
                    ${f.name}
                </p>
                <p>
                    Valid until: ${fileDies}
                </p>
                <div>
                    Compression: ${f.compression}
                </div>
                </div>
            </div>
            `
            )
        });
    })
        .catch((err) => {
            console.log("Something went horribly wrong: " + err);
        })
}).resolve();

router.on('/', () => {
    document.getElementById('file_list').innerHTML = '';
}).resolve();
