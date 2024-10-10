if ('VKIDSDK' in window) {
    const VKID = window.VKIDSDK;

    VKID.Config.init({
        app: 52449762,
        redirectUrl: 'http://localhost/register',
        responseMode: VKID.ConfigResponseMode.Callback,
        source: VKID.ConfigSource.LOWCODE,
    });

    const oneTap = new VKID.OneTap();

    oneTap.render({
        container: document.currentScript.parentElement,
        showAlternativeLogin: true,
        scheme: 'dark',
    })
    .on(VKID.WidgetEvents.ERROR, vkidOnError)
    .on(VKID.OneTapInternalEvents.LOGIN_SUCCESS, function (payload) {
        const code = payload.code;
        const deviceId = payload.device_id;

        VKID.Auth.exchangeCode(code, deviceId)
        .then(vkidOnSuccess)
        .catch(vkidOnError);
    });
}

function vkidOnSuccess(data) {
    fetch('/register', {
        method: 'POST',
        body: JSON.stringify(data)
    })
    .then(response => {
        if (response.ok) {
            window.location.href = '/folder';
        } else {
            console.error('Ошибка регистрации:', response.statusText);
        }
    })
}

function vkidOnError(error) {
    console.log(error)
}