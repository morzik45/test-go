var app = new Vue({
  el: "#app",
  data() {
    return {
      listVariants: {},
      showList: true,
      current_variant_id = "",
      current_task_id = "",
    };
  },
  mounted() {
    this.getListVariants();
  },
  methods: {
    singout() {
      axios
        .put("/api")
        .then((response) => {
          window.location.href = '/'
        })
        .catch((error) => {
          console.log(error);
        });
    },
    getListVariants() {
      axios
        .get("/api")
        .then((response) => {
          this.listVariants = response.data.variants;
        })
        .catch((error) => {
          console.log(error);
        });
    },
    getTask(variant_id, task_id) {
      axios
        .get("/api", {
          params: {
            variant_id: variant_id,
            task_id: task_id,
          },
        })
        .then((response) => {
          console.log(response.data);
        })
        .catch(function (error) {
          console.log(error);
        });
    },
  },
});
