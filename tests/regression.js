const phantom = require('phantom');
const sleep = require('sleep');
const chai = require('chai');
const assert = chai.assert;
const request = require('request-promise');
const colors = require('colors');

//Run the tests in an async issolated function so that we can use await semantics for promise results
(async function() {

    console.log("Starting basic functionality test !".cyan);

    //Create an instance of phantom and open a page
    const instance = await phantom.create();
    const page = await instance.createPage();

    //Open and authenticate, assumes the admin accounts is 'admin' and password 'admin'
    const status_auth = await page.open('http://localhost:3280/get/authentication/?identifier=admin&credentials=admin');
    const content = await page.property('content');
    await page.property('viewportSize', {width: 1600, height: 1240})

    const title = await page.evaluate(function() { return document.getElementsByClassName('title_header')[0].textContent } );

    assert.typeOf(title, 'string');
    assert.equal(title.toLowerCase(), 'hermes');
    console.log('We are on the right page'.green.dim);

    const auth_cookie = await page.evaluate(function() { return document.cookie } );
    assert.equal(auth_cookie.indexOf('auth=admin#|#'), 0);
    console.log("We have the authentication cookie".green.dim);

    const status_upload = await page.uploadFile('#file', __dirname + '/test_upload.txt');
    await page.evaluate(function() {
        document.querySelector('#upload_form > div:nth-child(1) > input[type="text"]').value = 'testFile.txt';
        document.querySelector('#upload_form > div.fields > div:nth-child(2) > input[type="text"]').value = '2400';
        document.querySelector('#public_switch').click();
        document.querySelector('#submit_form').click();
    });

    await instance.exit();

    //Phantom doesn't support my javascript so now its time for api testing using the browser
    //Also sleep seems to break if I call it here :/
    setTimeout(async () => {

        const file_list_api_answer = await request('http://localhost:3280/get/list/')

        assert.typeOf(file_list_api_answer, 'string');
        assert.equal(file_list_api_answer.indexOf('testFile.txt|#|'), 0);
        console.log("The file list seems to be in order".green.dim);

        const file_content_answer = await request('http://localhost:3280/get/file/?file=testFile.txt');

        assert.equal(file_content_answer, 'This is a test upload file\n');
        console.log("The file seems to be in order".green.dim);

        console.log("The basic functionality test has finished with sucessfully !".cyan);
    }, 2000);

}()).catch((e) => {
    console.log(e);
    process.exit();
});
