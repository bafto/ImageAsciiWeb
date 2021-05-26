let UploadButton = document.getElementById("uploadButton")
let output = document.getElementById("output")

UploadButton.addEventListener("click", async () => {
    let name = document.getElementById("fileInput").value
    let img = document.getElementById("fileInput").files[0]
    if (name != "") {
        let formData = new FormData()
        formData.append("image", img)
        let resp = await fetch("/image", {
            method: "POST",
            body: formData
        })
        let text = await resp.text()
        output.innerText = text
    }
})