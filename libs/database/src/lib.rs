use crate::generated::queries::query::{get_submission_by_id, Submission};
use deadpool_postgres::Pool;
use uuid::Uuid;

pub async fn fetch_submission(pool: &Pool, id: Uuid) -> Result<Option<Submission>, Error> {
    let client = pool.get().await?;
    
    let submission = get_submission_by_id()
        .bind(&client, &id)
        .opt()  // use .one() if always exists, .all() for Vec
        .await?;
    
    Ok(submission)
}