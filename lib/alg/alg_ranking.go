package alg

import (
	"math"
)

// 基于 Hacker GetWikis 的排序算法
// 基于：点赞数、时间排序，本身就加入了时间，所以不用再联合 created_at 排序了
// 缺点：需要周期（每小时或每天）对全部文章/评论等更新评分
// https://www.ruanyifeng.com/blog/2012/02/ranking_algorithm_hacker_news.html
func SimpleRank(likes uint, ) {

}

// Wilson Interval 威尔逊区间排序算法
// 优点：根据正反向反馈（喜欢和不喜欢）来评分排序，并且引入标准方差，保证评论数越多评分越高，不依赖时间，每次有人反馈（喜欢或不喜欢）重新计算一次即可。
// 缺点：存在冷启动问题
//  联合 order by rank DESC, wid(or created_at) DESC 解决冷启动，按时间排序的问题
// Wilson Confidence Interval considers binomial distribution for score calculation i.e. it considers only positive and negative ratings. If your product is rated on a 5 scale rating, then we can convert ratings {1–3} into negative and {4,5} to positive rating and can calculate Wilson score.
// https://www.evanmiller.org/how-not-to-sort-by-average-rating.html
// https://medium.com/tech-that-works/wilson-lower-bound-score-and-bayesian-approximation-for-k-star-scale-rating-to-rate-products-c67ec6e30060
// @param positive 点赞的人数
// @param total 总反馈人数（正反馈+负反馈） total number of ratings
// @confidence 默认为 95% 置信度
// @return   (-1,1)   ==> 约等于  upvotes / total，

func LowerBound(upvotes, total uint, confidence float64) float64 {
	if total == 0 || total < upvotes {
		return 0.0
	}
	n := float64(total)

	// 平均值有3个标准差之内的范围。称为“68-95-99.7法则”或“经验法则”  95% 是经验值
	z := 1.96 // quantile of the standard normal distrution, confidence=0.95 经验值
	if confidence != 0.0 && confidence != 0.95 {
		z = StandardQuantile(confidence)
		//z = math.Erfinv(1.0 - (1.0-confidence)/2.0) // 这个好像有问题，0.95的结果不是 1.96
	}

	phat := float64(upvotes) / n
	x := z * z / (4.0 * n)
	score := (phat + 2.0*x - z*math.Sqrt((phat*(1.0-phat)+x)/n)) / (1.0 + 4.0*x)
	return score
}

// 采用  int，9位数字
func WilsonRanking(upvotes, total uint, confidence float64) int {
	lb := LowerBound(upvotes, total, confidence)
	rank := int(math.Round(lb * 1000000000.0))
	return rank
}
