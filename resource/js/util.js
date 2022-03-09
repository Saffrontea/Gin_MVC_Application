function textLimit(o) {
    if (o.innerText.length > 20){
        o.innerText = o.innerText.substring(0,20) + "..."
    }
}
addEventListener("load",()=>{
    document.querySelectorAll(".text-limit").forEach(function (e){
    textLimit(e)
})})
