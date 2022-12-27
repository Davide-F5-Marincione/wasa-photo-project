<script>
export default {
    data: function () {
        return {
            errormsg: null,
            file: ""
        };
    },
    methods: {
        showImg() {
            const [file] = document.getElementById("file-input").files
            this.file = file
            if (this.file) {
                var shower = document.getElementById("img-show")
                shower.src = URL.createObjectURL(this.file)
            }
        },
        async uploadImage() {
            if (!this.file) {
                return
            }
            title = document.getElementById("title").value
            if (!title) {
                return
            }

            var form = new FormData();
            form.append("photo", this.file)
            form.append("title", title)

            this.errormsg = null;
			try {
				await this.$axios.post("/photos", form, { headers: {"Content-Type": "multipart/form-data"}});
                this.$router.push("stream")
			} catch (e) {
				this.errormsg = e.toString();
			}
        }
    }
}
</script>

<template>
	<TopBar></TopBar>
	<div>
        <div class="upload">
            <input id="title" class="upload-title" type="text" placeholder="Insert title here" maxlength="64"/>
            <input id="file-input" accept="image/*" class="upload-photo" type="file" placeholder="Load photo here" :onchange="() => showImg()"/>
            <div class="upload-photo-container">
                <img id="img-show" class="photo" alt="Here should be your image..." src="#"/>
            </div>
            <button class="upload-button" :onclick="() => uploadImage()"> Upload!</button>
        </div>
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>