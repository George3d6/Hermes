const enableAuthenticationForm = () => {

    document.getElementById("submit_auth").addEventListener("click",(e) => {
        e.preventDefault();
        const identifier = document.getElementById("identifier_field").value;
        const credentials = document.getElementById("credentials_field").value;
        const xhr = new XMLHttpRequest();
        xhr.open("GET", `/get/authentication//?identifier=${identifier}&credentials=${credentials}`);
        xhr.send();
        xhr.onreadystatechange = () => {
            $('#sign_in_form_modal')
              .modal('hide')
            ;
        }
    });

    document.getElementById("submit_auth_make").addEventListener("click",(e) => {
        e.preventDefault();
        const identifier = document.getElementById("identifier_field_make").value;
        const credentials = document.getElementById("credentials_field_make").value;
        const uploadNumber = document.getElementById("uploadNumber_field_make").value;
        const uploadSize = document.getElementById("uploadSize_field_make").value;
        const reader = document.getElementById("reader_field_make").checked;
        const writer = document.getElementById("writer_field_make").checked;
        const admin = document.getElementById("admin_field_make").checked;
        const xhr = new XMLHttpRequest();
        xhr.open("GET", `/post/token/?identifier=${identifier}&credentials=${credentials}&uploadNumber=${uploadNumber}`
        + `&uploadSize=${uploadSize}&reader=${reader}&writer=${writer}&admin=${admin}`);
        xhr.send();
        xhr.onreadystatechange = () => {
            $('#permission_view')
              .modal('hide')
            ;
        }
    });

    document.getElementById("close_permission_view_holder").addEventListener("click", (e) => {
        $('#permission_view')
          .modal('hide')
        ;
    })
}


export default enableAuthenticationForm;
