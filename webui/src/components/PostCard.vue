<script>
import ImageContainer from "./ImageContainer.vue";
export default {
    props: ["imgId"],
    data: function () {
        return {
            errormsg: null,
			title:"Wait... is something missing?",
			author:"The browser itself",
			date: new Date().toLocaleString( 'sv', { timeZoneName: 'short' } ).split(" ").slice(0,2).join(" "),
			liked: false,
        };
    },
    methods: {
		async refreshData() {
			this.errormsg = null;
			try {
				var response = await this.$axios.get("/photos/" + this.imgId.toString());
				this.title = response.data["photo-title"]
				this.author = response.data["photo-author"]
				var auLink = document.getElementById("author-link")
				auLink.onclick = this.moveToAuthor;
				this.date = this.convDate(response.data["photo-date"])
				this.liked = response.data["liked"]
				var but = document.getElementById("button-like")
				if (this.liked) {
					but.textContent = "unlike"
					but.classList.add("post-unlike-button")
					but.classList.remove("post-like-button")
					but.onclick = () => this.unlikeThis().then()
				} else {
					but.textContent = "like"
					but.classList.remove("post-unlike-button")
					but.classList.add("post-like-button")
					but.onclick = () => this.likeThis().then()
				}
			} catch (e) {
				this.errormsg = e.toString();
			}
		},
		async likeThis() {
			try {
				await this.$axios.put("/photos/" + this.imgId.toString() + "/likes/" + localStorage.username);
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.refreshData();
			this.disableLike();
		},
		async unlikeThis() {
			try {
				await this.$axios.delete("/photos/" + this.imgId.toString() + "/likes/" + localStorage.username);
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.refreshData();
			this.disableLike();
		},
		disableLike() {
			var but = document.getElementById("button-like")
			but.disabled = true;
			setTimeout(function() {
				but.disabled = false;
			}, 500);
		},
		convDate(date) {
			return date //TODO: convert to local timezone
		},
		moveToAuthor() {
			this.$router.push({name: 'user', params: {username: this.author}})
		}
    },
	created() {
		if (this.imgId) {
			this.refreshData();
		}
	},
    components: { ImageContainer }
}
</script>


<template>
	<div class="post-card">
		<div class="post-titleandinfo">
			<div class="post-title">
				{{ title }}
			</div>
			<div class="post-info">
				<div id="author-link" class="post-author" style="cursor: pointer;">
					{{ author }}
				</div>
				<div class="post-dateandlike">
					<div class="post-date">
						{{ date }}
					</div>
					<button id="button-like" class="post-like-button" liked="">
						Ahem...
					</button>
				</div>
			</div>
		</div>
		<div class="post-photoandcomments">
			<ImageContainer :alt="title" :img-id="imgId"></ImageContainer>
			<div class="post-comments">
				{{ resp }}
			</div>
			<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
		</div>
	</div>
</template>