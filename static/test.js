var app = new Vue({
  el: "#app",
  data() {
    return {
      listVariants: {},
    };
  },
  mounted() {
    this.getListVariants();
  },
  methods: {
    getListVariants() {
      axios
        .get("/")
        .then((response) => {
          console.log(this.listVariants);
          this.listVariants = response.data;
          console.log(this.listVariants);
        })
        .catch((error) => {
          console.log(error);
        });
    },
    getTask(variant_id, task_id) {
      axios
        .get("/", {
          params: {
            variant: variant_id,
            task: task_id,
          },
        })
        .then(function (response) {
          console.log(response.data);
          this.message = response.data;
        })
        .catch(function (error) {
          console.log(error);
        });
    },
  },
});
