const httpRoot = "";
const DownloadFileName = "scp294-output.txt";

document.addEventListener("DOMContentLoaded", () => {
    onGroupChange();
    onInputTypeChange();
    let inputType = document.getElementById("inputType");
    window.fileType = inputType.options[inputType.options.length - 1].value;
});

function convert() {
    let inputType = document.getElementById("inputType");
    let inputFormat = document.getElementById("inputFormat");
    let outputType = document.getElementById("outputType");
    let outputFormat = document.getElementById("outputFormat");
    let inputText = document.getElementById("inputText");
    let formData = new FormData();
    formData.append("InputType", inputType.value);
    formData.append("InputFormat", inputFormat.value);
    formData.append("OutputType", outputType.value);
    formData.append("OutputFormat", outputFormat.value);
    let inputIsFile = inputType.value == window.fileType;
    let outputIsFile = outputType.value == window.fileType;
    if (inputIsFile) {
        let inputFile = document.getElementById("inputFile");
        let files = inputFile.files;
        if (files != null && files.length > 0) {
            formData.append("InputFile", files[0]);
        } else {
            alert("No file to convert")
            return;
        }
    } else {
        if (inputText.value != null && inputText.value != "") {
            formData.append("InputData", inputText.value);
        } else {
            alert("Nothing to convert")
            return;
        }
    }

    fetch(httpRoot + "/convert", {
        method: "POST",
        body: formData,
    }).then(re => {
        if (re.ok) {
            if (inputIsFile) {
                if (outputIsFile) {
                    return streamToFile(re);
                } else {
                    return streamToText(re);
                }
            } else {
                if (outputIsFile) {
                    return textToFile(re);
                } else {
                    return re.json();
                }
            }
        }

        return re.json();
    }).then(re => {
        if (re.Success) {
            setOutput(re.Data);
        }
        updateMessage(re)
    })
}

async function readStream(re) {
    let success = true;
    let errorMessage = null;

    let reader = re.body.getReader();
    let myTypedArray = new MyTypedArray(Uint8Array, 4096);
    let transferSize = 0;
    let done = false;
    while (!done) {
        done = await reader.read().then(({done, value}) => {
            if (value != null && value.byteLength > 0) {
                myTypedArray.pushArray(value);
                transferSize += value.byteLength;
                updateMessageValue("Transferred size: " + (transferSize >>> 10) + "KB", "blue");
            }
            if (done) {
                updateMessageValue("Transfer done, total size: " + (transferSize >>> 10) + "KB", "green");
                return true;
            }
            return false;
        }).catch(error => {
            success = false;
            errorMessage = error.message;
        });
    }

    return {
        Success: success,
        Message: success ? "Data was converted! Total size: " + (transferSize >>> 10) + "KB" : errorMessage,
        TypedArray: myTypedArray
    };
}

async function streamToFile(re) {
    let result = await readStream(re);
    let file = new File([result.TypedArray.toString()], DownloadFileName, {type: "text/plain"});
    let url = URL.createObjectURL(file);
    let a = document.createElement('a');
    a.href = url;
    a.download = DownloadFileName;
    a.click();
    URL.revokeObjectURL(url);

    return {
        Success: result.Success,
        Message: result.Message,
        Data: result.TypedArray.preview()
    };
}

async function streamToText(re) {
    let result = await readStream(re);

    return {
        Success: result.Success,
        Message: result.Message,
        Data: result.TypedArray.toString()
    };
}

async function textToFile(re) {
    let result = await re.json();
    let file = new File([result.Data], DownloadFileName, {type: "text/plain"});
    let url = URL.createObjectURL(file);
    let a = document.createElement('a');
    a.href = url;
    a.download = DownloadFileName;
    a.click();
    URL.revokeObjectURL(url);

    return result;
}

class MyTypedArray {
    constructor(typedArrayClass, cap) {
        this.typedArrayClass = typedArrayClass;
        this.off = 0;
        this.cap = cap;
        this.arrayBuffer = new ArrayBuffer(this.cap);
        this.typedArray = new this.typedArrayClass(this.arrayBuffer);
    }

    pushArray(valArr) {
        for (let val of valArr) {
            this.push(val);
        }
    }

    push(val) {
        if (this.off < this.cap) {
            this.typedArray[this.off++] = val;
            return;
        }
        this.cap *= 2;
        let newArrayBuffer = new ArrayBuffer(this.cap);
        let newTypedArray = new this.typedArrayClass(newArrayBuffer);
        newTypedArray.set(this.typedArray, 0);
        this.arrayBuffer = newArrayBuffer;
        this.typedArray = newTypedArray;
        this.push(val);
    }

    preview(charset) {
        if (charset == null) {
            charset = "utf-8";
        }
        let decoder = new TextDecoder(charset, {ignoreBOM: true})
        return "Preview the first 65536 bytes:\n" + decoder.decode(new this.typedArrayClass(this.arrayBuffer, 0, 65536));
    }

    toString(charset) {
        if (charset == null) {
            charset = "utf-8";
        }
        let decoder = new TextDecoder(charset, {ignoreBOM: true})
        return decoder.decode(new this.typedArrayClass(this.arrayBuffer, 0, this.off));
    }
}

function clearText() {
    let inputText = document.getElementById("inputText");
    let inputFile = document.getElementById("inputFile");
    let output = document.getElementById("output");
    inputText.value = null;
    inputFile.value = null;
    output.value = null;
    updateMessageValue(null);
}

function loadRecord() {
    let selectRecord = document.getElementById("selectRecord");
    if (selectRecord.value == null || selectRecord.value == "") {
        alert("Record id cannot be empty")
        return;
    }
    fetch(httpRoot + "/loadRecord?RecordId=" + selectRecord.value, {
        method: "GET",
    }).then(re => {
        if (re.ok) return re.json();
    }).then(re => {
        if (re.Success) {
            let inputType = document.getElementById("inputType");
            let outputType = document.getElementById("outputType");
            inputType.value = re.Data.InputType;
            outputType.value = re.Data.OutputType;
            if (inputType.value != window.fileType) {
                let inputText = document.getElementById("inputText");
                inputText.value = re.Data.InputData;
            }
            onInputTypeChange();
            setOutput(re.Data.OutputData);
        }
        updateMessage(re)
    })
}

function addRecord() {
    let inputType = document.getElementById("inputType");
    if (inputType.value == window.fileType) {
        alert("Cannot save file record");
        return;
    }
    let recordName = prompt("Input a name", "");
    if (recordName == null || recordName == "") {
        alert("The name cannot be empty");
        return;
    }
    let selectGroup = document.getElementById("selectGroup");
    let outputType = document.getElementById("outputType");
    let inputText = document.getElementById("inputText");
    let output = document.getElementById("output");
    let record = {
        Id: 0,
        Name: recordName,
        InputType: parseInt(inputType.value),
        OutputType: parseInt(outputType.value),
        InputData: inputText.value,
        OutputData: output.value,
        GroupId: Number(selectGroup.value),
    }
    let groupId = document.getElementById("selectGroup").value;
    fetch(httpRoot + "/addRecord", {
        method: "POST",
        body: JSON.stringify(record),
        headers: {"Content-Type": "application/json;charset=UTF-8",},
    }).then(re => {
        if (re.ok) return re.text();
    }).then(re => {
        let divHeader = document.getElementById("divHeader");
        divHeader.innerHTML = re;
        switchGroup(groupId);
        updateMessageValue("Record was saved!")
    })
}

function deleteRecord() {
    if (!confirm("Delete it?")) {
        return;
    }
    let selectRecord = document.getElementById("selectRecord");
    if (selectRecord.value == null || selectRecord.value == "") {
        alert("Record id cannot be empty")
        return;
    }
    let groupId = document.getElementById("selectGroup").value;
    fetch(httpRoot + "/deleteRecord?RecordId=" + selectRecord.value, {
        method: "DELETE",
    }).then(re => {
        if (re.ok) return re.text();
    }).then(re => {
        let divHeader = document.getElementById("divHeader");
        divHeader.innerHTML = re;
        switchGroup(groupId);
        updateMessageValue("Record was deleted!")
    })
}

function createGroup() {
    let groupName = prompt("Input a name", "");
    if (groupName == null || groupName == "") {
        alert("The name cannot be empty")
        return;
    }
    let group = {
        Id: 0,
        Name: groupName,
    }
    fetch(httpRoot + "/addGroup", {
        method: "POST",
        body: JSON.stringify(group),
        headers: {"Content-Type": "application/json;charset=UTF-8",},
    }).then(re => {
        if (re.ok) return re.text();
    }).then(re => {
        let divHeader = document.getElementById("divHeader");
        divHeader.innerHTML = re;
        updateMessageValue("Group was added!")
    })
}

function deleteGroup() {
    if (!confirm("Delete it?")) {
        return;
    }
    let selectGroup = document.getElementById("selectGroup");
    if (selectGroup.value == null || selectGroup.value == "") {
        alert("Group id cannot be empty")
        return;
    }
    fetch(httpRoot + "/deleteGroup?GroupId=" + selectGroup.value, {
        method: "DELETE",
    }).then(re => {
        if (re.ok) return re.text();
    }).then(re => {
        let divHeader = document.getElementById("divHeader");
        divHeader.innerHTML = re;
        updateMessageValue("Group was deleted!")
    })
}

function switchGroup(groupId) {
    let selectGroup = document.getElementById("selectGroup");
    for (let option of selectGroup.options) {
        if (option.value == groupId) {
            option.selected = true;
            onGroupChange();
            break;
        }
    }
}

function onGroupChange() {
    let selectGroup = document.getElementById("selectGroup");
    let selectRecord = document.getElementById("selectRecord");
    let selected = false;
    for (let option of selectRecord.options) {
        option.selected = false;
        if (option.getAttribute("groupId") == selectGroup.value) {
            option.style.display = null;
            if (!selected) {
                option.selected = true;
                selected = true;
            }
        } else {
            option.style.display = "none";
        }
    }
    if (!selected) {
        selectRecord.selectedIndex = -1;
    }
}

function onInputTypeChange() {
    let inputType = document.getElementById("inputType");
    let inputText = document.getElementById("inputText");
    let inputFileDiv = document.getElementById("inputFileDiv");
    if (inputType.value != window.fileType) {
        inputText.style.display = null;
        inputFileDiv.style.display = "none";
    } else {
        inputText.style.display = "none";
        inputFileDiv.style.display = null;
    }
}

function setOutput(data) {
    let output = document.getElementById("output");
    output.value = data;
    if (output.scrollHeight > 200) {
        output.style.height = '200px';
        output.style.height = output.scrollHeight + 100 + 'px';
    }
}

function updateMessage(re) {
    let element = document.getElementById("message");
    element.innerText = re.Message;
    element.style.color = re.Success ? "green" : "red";
}

function updateMessageValue(value, color) {
    let element = document.getElementById("message");
    element.innerText = value;
    if (color != null) {
        element.style.color = color;
    } else if (value != null) {
        element.style.color = "green";
    }
}