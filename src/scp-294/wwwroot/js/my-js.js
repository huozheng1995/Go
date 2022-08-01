const httpRoot = "";

document.addEventListener("DOMContentLoaded", () => {
    onGroupChange();
    onInputTypeChange();
});

function convert() {
    let inputType = document.getElementById("inputType");
    let outputType = document.getElementById("outputType");
    let input = document.getElementById("input");
    let formData = new FormData();
    formData.append("InputType", inputType.value);
    formData.append("OutputType", outputType.value);
    let isFile = inputType.value == "File";
    if (!isFile) {
        if (input.value != null && input.value != "") {
            formData.append("InputData", input.value);
        } else {
            alert("Nothing to convert")
            return;
        }
    } else {
        let inputFile = document.getElementById("inputFile");
        let files = inputFile.files;
        if (files != null && files.length > 0) {
            formData.append("InputFile", files[0]);
        } else {
            alert("No file to convert")
            return;
        }
    }

    fetch(httpRoot + "/convert", {
        method: "POST",
        body: formData,
    }).then(re => {
        if (!isFile) {
            if (re.ok) {
                return re.json();
            }
        } else {
            if (re.ok) {
                return readTextStream(re);
            } else {
                return re.json();
            }
        }
    }).then(re => {
        if (re.Success) {
            setOutput(re.Data);
        }
        updateMessage(re)
    })
}

async function readTextStream(re) {
    let reader = re.body.getReader();
    let myTypedArray = new MyTypedArray(Uint8Array, 4096);
    let transferSize = 0;
    let done = false;
    while (!done) {
        done = await reader.read().then(result => {
            if (result.value != null && result.value.byteLength > 0) {
                myTypedArray.pushArray(result.value);
                transferSize += result.value.byteLength;
                updateMessageValue("Transferred size: " + (transferSize >>> 10) + "KB", "blue");
            }
            if (result.done) {
                updateMessageValue("Transfer done, total size: " + (transferSize >>> 10) + "KB", "green");
                return true;
            }
            return false;
        });
    }
    return {
        Success: true,
        Message: "Data was converted! Total size: " + (transferSize >>> 10) + "KB",
        Data: myTypedArray.toString()
    };
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

    toString(charset) {
        if (charset == null) {
            charset = "utf-8";
        }
        let decoder = new TextDecoder(charset, {ignoreBOM: true})
        return decoder.decode(new this.typedArrayClass(this.arrayBuffer, 0, this.off));
    }
}

function clearText() {
    let input = document.getElementById("input");
    let inputFile = document.getElementById("inputFile");
    let output = document.getElementById("output");
    input.value = null;
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
            if (inputType.value != "File") {
                let input = document.getElementById("input");
                input.value = re.Data.InputData;
            }
            onInputTypeChange();
            setOutput(re.Data.OutputData);
        }
        updateMessage(re)
    })
}

function addRecord() {
    let inputType = document.getElementById("inputType");
    if (inputType.value == "File") {
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
    let input = document.getElementById("input");
    let output = document.getElementById("output");
    let record = {
        Id: 0,
        Name: recordName,
        InputType: inputType.value,
        OutputType: outputType.value,
        InputData: input.value,
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
    let input = document.getElementById("input");
    let inputFileDiv = document.getElementById("inputFileDiv");
    if (inputType.value != "File") {
        input.style.display = null;
        inputFileDiv.style.display = "none";
    } else {
        input.style.display = "none";
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