export const displayError = (message: string) => {
    document.getElementById('error_holder').innerHTML = '';
    document.getElementById('error_holder').insertAdjacentHTML('afterbegin', `
    <div class="ui error message">
    <i class="close icon" id="close_message"></i>
        <div class="header">
            There was an error
        </div>
        <p>
            ${message}
        </p>
    </div>
    `);
    setTimeout(() => {
        document.getElementById('error_holder').innerHTML = '';
    }, 2600);
}

export const displayMessage = (message: string) => {
    document.getElementById('error_holder').innerHTML = '';
    document.getElementById('error_holder').insertAdjacentHTML('afterbegin', `
    <div class="ui info message">
    <i class="close icon" id="close_message"></i>
        <div class="header">
            Success
        </div>
        <p>
            ${message}
        </p>
    </div>
    `);
    setTimeout(() => {
        document.getElementById('error_holder').innerHTML = '';
    }, 2600);
}
