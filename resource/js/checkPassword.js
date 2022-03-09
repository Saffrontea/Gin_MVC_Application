
var elm_pass = $('#password');
var elm_confirm = $('#confirm-password');
var invalid_feedback = $('.invalid-feedback')[0];
const invalidText = invalid_feedback.innerText;
const charLengthInvalidText = "英数字8文字以上で入力してください。";

/*
* 確認パスワードのキーボード（KeyUp）イベントリスナー
*/


elm_confirm.on('keyup', function () {
    // まだパスワード（確認）を入力していない
    if (elm_confirm.val() === "") {
        elm_confirm.removeClass("is-valid");
        elm_confirm.removeClass("is-invalid");
        return;
    }
    // 先頭から一文字ずつ取り出してチェックし最後まで到達していなくとも「問題無し」と判断
    var array_pass_chars = elm_pass.val().split("");
    var array_confirm_chars = elm_confirm.val().split("");
    $.each(array_confirm_chars, function (index, char) {
        if (array_pass_chars[index] === char) {
            // 先頭から一文字ずつ一致している場合には中途でも何も表示しない
            elm_confirm.removeClass("is-valid");
            elm_confirm.removeClass("is-invalid");
        } else {
            // 一文字でも異なる場合はInvalid
            elm_confirm.removeClass("is-valid");
            invalid_feedback.innerText = invalidText;
            elm_confirm.addClass("is-invalid");
            return false;
        }
    });
    // 完全一致したらValid
    if (elm_pass.val() === elm_confirm.val()) {
        if (elm_pass.val().length >= 8){ //length足りる?
            elm_confirm.addClass("is-valid");
        }else{
            invalid_feedback.innerText = charLengthInvalidText;
            elm_confirm.addClass("is-invalid");
        }
    } else {
        invalid_feedback.innerText = invalidText;
        elm_confirm.addClass("is-invalid");
    }
});

elm_pass.on('keyup', function () {
    // まだパスワード（確認）を入力していない
    if (elm_confirm.val() === "") {
        elm_confirm.removeClass("is-valid");
        elm_confirm.removeClass("is-invalid");
        return;
    }
    // 先頭から一文字ずつ取り出してチェックし最後まで到達していなくとも「問題無し」と判断
    var array_pass_chars = elm_pass.val().split("");
    var array_confirm_chars = elm_confirm.val().split("");
    $.each(array_confirm_chars, function (index, char) {
        if (array_pass_chars[index] === char) {
            // 先頭から一文字ずつ一致している場合には中途でも何も表示しない
            elm_confirm.removeClass("is-valid");
            elm_confirm.removeClass("is-invalid");
        } else {
            // 一文字でも異なる場合はInvalid
            elm_confirm.removeClass("is-valid");
            invalid_feedback.innerText = invalidText;
            elm_confirm.addClass("is-invalid");
            return false;
        }
    });
    // 完全一致したらValid
    if (elm_pass.val() === elm_confirm.val()) {
        if (elm_pass.val().length >= 8){ //length足りる?
            elm_confirm.addClass("is-valid");
        }else{
            invalid_feedback.innerText = charLengthInvalidText;
            elm_confirm.addClass("is-invalid");
        }
    } else {
        invalid_feedback.innerText = invalidText;
        elm_confirm.addClass("is-invalid");
    }
});

/*
* 確認パスワード入力のフォーカスを失ったとき（Blur）のイベントリスナー
*/
elm_confirm.on('blur', function () {
    if (elm_pass.val() === elm_confirm.val()) {
        elm_confirm.removeClass("is-invalid");
        if (elm_pass.val().length >= 8){ //length足りる?
            elm_confirm.addClass("is-valid");
        }else{
            invalid_feedback.innerText = charLengthInvalidText;
            elm_confirm.addClass("is-invalid");
        }
    }
    // フォーカスが失われたときパスワードが一致しないと判断（Invalid）
    else {
        elm_confirm.removeClass("is-valid");
        invalid_feedback.innerText = invalidText;
        elm_confirm.addClass("is-invalid");
    }
});
