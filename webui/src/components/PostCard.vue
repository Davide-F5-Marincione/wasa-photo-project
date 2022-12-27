<script>
import ImageContainer from "./ImageContainer.vue";
export default {
    props: ["imgId", "del"],
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
			commentsLimit: "",
			username: localStorage.username,
			commentWrite: ""
        };
    },
    methods: {
		async refreshData(new_likes=false, new_comments=false) {
			if (!this.imgId) {
				return
			}
			this.errormsg = null;
			try {
				var response = await this.$axios.get("/photos/" + this.imgId.toString(), {params: {"likes-base": this.likesBase, "comments-limit":this.commentsLimit}});
				this.title = response.data["photo-title"]
				this.author = response.data["photo-author"]
				this.date = this.convDate(response.data["photo-date"])
				this.liked = response.data["liked"]
				
				for (let i=0; i< response.data["likes-running-batch"].length; i++) {
					if (response.data["likes-running-batch"][i].name == this.username) {
						response.data["likes-running-batch"].splice(i, 1)
						break
					}
				}

				if (new_likes) {
					this.likersResult.push(...response.data["likes-running-batch"])
					if (response.data["likes-running-batch"].length > 0) {
						this.likesBase = response.data["likes-running-batch"][response.data["likes-running-batch"].length - 1].name;
					}
				}
				if (new_comments) {
					this.commentsResults.push(...response.data["comments-running-batch"])
					if (response.data["comments-running-batch"].length > 0) {
						this.commentsLimit = response.data["comments-running-batch"][response.data["comments-running-batch"].length - 1]["comment-id"];
					}
				}

			} catch (e) {
				this.errormsg = e.toString();
			}
		},
		async likeThis() {
			if (!this.imgId) {
				return
			}
			try {
				await this.$axios.put("/photos/" + this.imgId.toString() + "/likes/" + localStorage.username);
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.refreshData();
		},
		async unlikeThis() {
			if (!this.imgId) {
				return
			}
			try {
				await this.$axios.delete("/photos/" + this.imgId.toString() + "/likes/" + localStorage.username);
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.refreshData();
		},
		convDate(date) {
			return (new Date(date.split(" ").join("T") + "Z")).toLocaleString( 'sv', { timeZoneName: 'short' } ).split(" ").slice(0,2).join(" ") //TODO: convert to local timezone
		},
		async sendComment() {
			if (!this.imgId) {
				return
			}
			var textarea = document.getElementById("comment-text");
			try {
				await this.$axios.post("/photos/" + this.imgId.toString() + "/comments",  textarea.value, {headers: {'Content-Type': 'application/json'}});
				textarea.value = ""
				this.commentsResults = []
				this.commentsLimit = ""
				this.refreshData(false, true)
			} catch (e) {
				this.errormsg = e.toString();
			}
		},

		async deleteComment(id) {
			if (!this.imgId) {
				return
			}

			try {
				await this.$axios.delete("/photos/" + this.imgId.toString() + "/comments/" + id.toString());
				this.commentsResults = []
				this.commentsLimit = ""
				this.refreshData(false, true)
			} catch (e) {
				this.errormsg = e.toString();
			}
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
						<router-link class="post-author" align="right" :to="'/users/' + author">{{ author }}</router-link>
						<div class="post-date" align="right">
							{{ date }}
						</div>
						<button v-if="author == username" class="delete-post" align="right" :onclick="() => del()">delete this post</button>
					</div>
				</div>
			</div>
			<ImageContainer :alt="title" :img-id="imgId"></ImageContainer>
		</div>
		<div class="post-commentsandlikes">
			<div id="likesholder" class="post-likesholder">
				<router-link class="like-element" v-for="element in likersResult"  v-bind="element.name" :to="'/users/' +element.name">{{ element.name }}</router-link>
				<router-link v-if="liked" class="like-element" :to="'/users/' + username">{{ username }}</router-link>
				<button class="like-element-end" :onclick="() => this.refreshData(true, false)">More likes!</button>
			</div>
			<button v-if="!liked" class="post-like-button" :onclick="() => this.likeThis().then()">like</button>
			<button v-if="liked" class="post-unlike-button" :onclick="() => this.unlikeThis().then()">unlike</button>
			<div class="post-comments disable-scrollbars">
				<div v-for="element in commentsResults" v-bind="element['comment-id']" class="post-comment">
					<div class="comment-authoranddeleter">
						<router-link class="comment-author" :to="'/users/' + element['comment-author']"> {{ element["comment-author"] }}: </router-link>
						<button v-if="element['comment-author']==username" class="comment-deleter" :onclick="() => this.deleteComment(element['comment-id'])">✕</button>
					</div>
					<div class="comment-content">{{ element["comment-text"] }}</div>
					<div class="comment-date">{{ convDate(element["comment-date"]) }}</div>
				</div>
				<button class="post-comments-end" :onclick="() => this.refreshData(false, true)">See more comments!</button>
			</div>
			<textarea id="comment-text" class="post-comment-writer" maxlength="256"></textarea>
			<button class="post-comment-button" :onclick="() => this.sendComment()">
				Comment!
			</button>
			<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
		</div>
	</div>
</template>