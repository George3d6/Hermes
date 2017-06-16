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
        const equal = document.getElementById("equal_field_make").value;
        const admin = document.getElementById("admin_field_make").value;
        const xhr = new XMLHttpRequest();
        xhr.open("GET", `/authenticate/?identifier=${identifier}&credentials=${credentials}&uploadNumber=${uploadNumber}`
        + `&uploadSize=${uploadSize}&equal=${equal}&admin=${admin}`);
        xhr.send();
        xhr.onreadystatechange = () => {
            document.getElementById('permission_view').display = "none";
        }
    });
}


export default enableAuthenticationForm;
