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
			likersResult: [],
			commentsResults: [],
			likesBase: "",
			commentsLimit: ""
        };
    },
    methods: {
		async refreshData(new_likes=false, new_comments=false) {
			this.errormsg = null;
			try {
				var response = await this.$axios.get("/photos/" + this.imgId.toString(), {params: {"likes-base": this.likesBase, "comments-limit":this.commentsLimit}});
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
				
				if (new_likes) {
					this.likersResult.push(...response.data["likes-running-batch"])
					this.likesBase = response.data["likes-running-batch"][response.data["likes-running-batch"].length - 1]["name"];
				}
				if (new_comments) {
					this.commentsResults.push(...response.data["comments-running-batch"])
					this.commentsLimit = response.data["comments-running-batch"][response.data["comments-running-batch"].length - 1]["comment-id"];
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
			this.refreshData(true, true);
		}
	},
    components: { ImageContainer }
}
</script>


<template>
	<div class="post-card">
		<div class="post-titleinfoandphoto">
			<div class="post-titleandinfo">
				<div class="post-title">
					{{ title }}
				</div>
				<div class="post-info">
					<div>					
						<div id="author-link" class="post-author" style="cursor: pointer;" align="right">
							{{ author }}
						</div>
						<div class="post-date" align="right">
							{{ date }}
						</div>
					</div>
				</div>
			</div>
			<ImageContainer :alt="title" :img-id="imgId"></ImageContainer>
		</div>
		<div class="post-commentsandlikes">
			<div id="likesholder" class="post-likesholder">
				<router-link class="like-element" v-for="element in likersResult" :to="'/users/' +element['name']">{{ element['name'] }}</router-link>
			</div>
			<button id="button-like" class="post-like-button">
				like
			</button> 
			<div class="post-comments">
			</div>
			<textarea class="post-comment-writer" maxlength="256"></textarea>
			<button id="button-comment" class="post-comment-button">
				Comment!
			</button> 
			<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
		</div>
	</div>
</template>

<!-- this.likersResult.forEach(element => {
	var liker = element["name"]
	var but = document.createElement("button");
	but.addEventListener('click', () => this.$router.push({name: 'user', params: {username: liker}}));
	var text = document.createTextNode(liker);
	but.classList.add("like-element")
	but.appendChild(text);
	holder.appendChild(but);
});

if (this.likersResult.length >= 64) {
	var but = document.createElement("button");
	but.addEventListener('click', () => this.refreshData(new_likes=true));
	but.setAttribute("id", "furtherLikes_button")
	var text = document.createTextNode("More likers!");
	but.classList.add("like-element-end")
	but.appendChild(text);
	holder.appendChild(but);
} -->