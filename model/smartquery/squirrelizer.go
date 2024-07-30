package smartquery

import (
	"github.com/Masterminds/squirrel"
)

type Squirrelizer struct {
	Sql string
}

func (r *Squirrelizer) BuildSelect(playlistID, userID, orderBy string) (squirrel.SelectBuilder, error) {
	if orderBy == "" {
		orderBy = "title"
	}
	// sq := Select("row_number() over (order by "+orderBy+") as id", "'"+pls.ID+"' as playlist_id", "media_file.id as media_file_id").
	// 	From("media_file").LeftJoin("annotation on (" +
	// 	"annotation.item_id = media_file.id" +
	// 	" AND annotation.item_type = 'media_file'" +
	// 	" AND annotation.user_id = '" + userId(r.ctx) + "')").
	// 	LeftJoin("media_file_genres ag on media_file.id = ag.media_file_id").
	// 	LeftJoin("genre on ag.genre_id = genre.id").GroupBy("media_file.id")
	sq := squirrel.Select("row_number() over (order by "+orderBy+") as id", "'"+playlistID+"' as playlist_id", "media_file.id as media_file_id").
		From("media_file").
		LeftJoin("annotation on (annotation.item_id = media_file.id AND annotation.item_type = 'media_file' " +
			"AND annotation.user_id = '" + userID + "')")
	// .
	// LeftJoin("media_file_genres ag on media_file.id = ag.media_file_id").
	// LeftJoin("genre on ag.genre_id = genre.id").GroupBy("media_file.id")
	sql, _, err := sq.ToSql()
	r.Sql = sql + " WHERE "
	return sq, err
}

func (r *Squirrelizer) BuildRefreshSmartQueryPlaylistSQL(playlistID, userID, smartQuery, orderBy string) (squirrel.InsertBuilder, error) {

	/* From the SQLite documentation
	 *		The row_number() window function assigns consecutive integers to each row in order of the "ORDER BY" clause within the window-defn
	 *       	SELECT x, y, row_number() OVER (ORDER BY y) AS row_number FROM t0 ORDER BY x;
	 * 			(in this case "ORDER BY y")
	 * 		Note that this does not affect the order in which results are returned from the overall query.
	 *  	The order of the final output is still governed by the ORDER BY clause attached to the SELECT statement (in this case "ORDER BY x").
	 */
	if orderBy == "" {
		orderBy = "title"
	}
	sq, err := r.BuildSelect(playlistID, userID, orderBy)
	if err != nil {
		return squirrel.InsertBuilder{}, err
	}

	sq = sq.Where(smartQuery)
	insSql := squirrel.Insert("playlist_tracks").Columns("id", "playlist_id", "media_file_id").Select(sq)
	sql, _, err := insSql.ToSql()

	// fmt.Println("**********************************")
	// fmt.Printf("(A) [%v]", smartQuery)
	// fmt.Println()
	// fmt.Println("**********************************")
	// fmt.Printf("(B) [%v]", sql)
	// fmt.Println()
	// fmt.Println("**********************************")

	r.Sql = sql
	return insSql, err
}
