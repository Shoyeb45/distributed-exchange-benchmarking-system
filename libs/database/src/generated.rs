// This file was generated with `cornucopia`. Do not modify.

#[allow(clippy::all, clippy::pedantic)] #[allow(unused_variables)]
#[allow(unused_imports)] #[allow(dead_code)] pub mod types { pub mod public { #[derive(serde::Serialize, Debug, Clone, Copy, PartialEq, Eq)]
#[allow(non_camel_case_types)] pub enum SupportedLanguage { CPP,RUST,}impl<'a> postgres_types::ToSql for SupportedLanguage
{
    fn
    to_sql(&self, ty: &postgres_types::Type, buf: &mut
    postgres_types::private::BytesMut,) -> Result<postgres_types::IsNull,
    Box<dyn std::error::Error + Sync + Send>,>
    {
        let s = match *self { SupportedLanguage::CPP => "CPP",SupportedLanguage::RUST => "RUST",};
        buf.extend_from_slice(s.as_bytes());
        std::result::Result::Ok(postgres_types::IsNull::No)
    } fn accepts(ty: &postgres_types::Type) -> bool
    {
        if ty.name() != "supported_language" { return false; } match *ty.kind()
        {
            postgres_types::Kind::Enum(ref variants) =>
            {
                if variants.len() != 2 { return false; }
                variants.iter().all(|v| match &**v
                { "CPP" => true,"RUST" => true,_ => false, })
            } _ => false,
        }
    } fn
    to_sql_checked(&self, ty: &postgres_types::Type, out: &mut
    postgres_types::private::BytesMut,) -> Result<postgres_types::IsNull,
    Box<dyn std::error::Error + Sync + Send>>
    { postgres_types::__to_sql_checked(self, ty, out) }
} impl<'a> postgres_types::FromSql<'a> for SupportedLanguage
{
    fn from_sql(ty: &postgres_types::Type, buf: &'a [u8],) ->
    Result<SupportedLanguage, Box<dyn std::error::Error + Sync + Send>,>
    {
        match std::str::from_utf8(buf)?
        {
            "CPP" => Ok(SupportedLanguage::CPP),"RUST" => Ok(SupportedLanguage::RUST),s =>
            Result::Err(Into::into(format!("invalid variant `{}`", s))),
        }
    } fn accepts(ty: &postgres_types::Type) -> bool
    {
        if ty.name() != "supported_language" { return false; } match *ty.kind()
        {
            postgres_types::Kind::Enum(ref variants) =>
            {
                if variants.len() != 2 { return false; }
                variants.iter().all(|v| match &**v
                { "CPP" => true,"RUST" => true,_ => false, })
            } _ => false,
        }
    }
}#[derive(serde::Serialize, Debug, Clone, Copy, PartialEq, Eq)]
#[allow(non_camel_case_types)] pub enum SubmissionStatus { UPLOADING,UPLOADED,BUILDING,BUILD_FAILED,BINARY_READY,SUCCESS,}impl<'a> postgres_types::ToSql for SubmissionStatus
{
    fn
    to_sql(&self, ty: &postgres_types::Type, buf: &mut
    postgres_types::private::BytesMut,) -> Result<postgres_types::IsNull,
    Box<dyn std::error::Error + Sync + Send>,>
    {
        let s = match *self { SubmissionStatus::UPLOADING => "UPLOADING",SubmissionStatus::UPLOADED => "UPLOADED",SubmissionStatus::BUILDING => "BUILDING",SubmissionStatus::BUILD_FAILED => "BUILD_FAILED",SubmissionStatus::BINARY_READY => "BINARY_READY",SubmissionStatus::SUCCESS => "SUCCESS",};
        buf.extend_from_slice(s.as_bytes());
        std::result::Result::Ok(postgres_types::IsNull::No)
    } fn accepts(ty: &postgres_types::Type) -> bool
    {
        if ty.name() != "submission_status" { return false; } match *ty.kind()
        {
            postgres_types::Kind::Enum(ref variants) =>
            {
                if variants.len() != 6 { return false; }
                variants.iter().all(|v| match &**v
                { "UPLOADING" => true,"UPLOADED" => true,"BUILDING" => true,"BUILD_FAILED" => true,"BINARY_READY" => true,"SUCCESS" => true,_ => false, })
            } _ => false,
        }
    } fn
    to_sql_checked(&self, ty: &postgres_types::Type, out: &mut
    postgres_types::private::BytesMut,) -> Result<postgres_types::IsNull,
    Box<dyn std::error::Error + Sync + Send>>
    { postgres_types::__to_sql_checked(self, ty, out) }
} impl<'a> postgres_types::FromSql<'a> for SubmissionStatus
{
    fn from_sql(ty: &postgres_types::Type, buf: &'a [u8],) ->
    Result<SubmissionStatus, Box<dyn std::error::Error + Sync + Send>,>
    {
        match std::str::from_utf8(buf)?
        {
            "UPLOADING" => Ok(SubmissionStatus::UPLOADING),"UPLOADED" => Ok(SubmissionStatus::UPLOADED),"BUILDING" => Ok(SubmissionStatus::BUILDING),"BUILD_FAILED" => Ok(SubmissionStatus::BUILD_FAILED),"BINARY_READY" => Ok(SubmissionStatus::BINARY_READY),"SUCCESS" => Ok(SubmissionStatus::SUCCESS),s =>
            Result::Err(Into::into(format!("invalid variant `{}`", s))),
        }
    } fn accepts(ty: &postgres_types::Type) -> bool
    {
        if ty.name() != "submission_status" { return false; } match *ty.kind()
        {
            postgres_types::Kind::Enum(ref variants) =>
            {
                if variants.len() != 6 { return false; }
                variants.iter().all(|v| match &**v
                { "UPLOADING" => true,"UPLOADED" => true,"BUILDING" => true,"BUILD_FAILED" => true,"BINARY_READY" => true,"SUCCESS" => true,_ => false, })
            } _ => false,
        }
    }
} }}#[allow(clippy::all, clippy::pedantic)] #[allow(unused_variables)]
#[allow(unused_imports)] #[allow(dead_code)] pub mod queries
{ pub mod query
{ use futures::{{StreamExt, TryStreamExt}};use futures; use cornucopia_async::GenericClient;#[derive(serde::Serialize, Debug, Clone, PartialEq,)] pub struct GetSubmissionById
{ pub id : i32,pub user_id : i32,pub language : super::super::types::public::SupportedLanguage,pub source_code : String,pub status : super::super::types::public::SubmissionStatus,pub created_at : time::OffsetDateTime,pub updated_at : time::OffsetDateTime,}pub struct GetSubmissionByIdBorrowed<'a> { pub id : i32,pub user_id : i32,pub language : super::super::types::public::SupportedLanguage,pub source_code : &'a str,pub status : super::super::types::public::SubmissionStatus,pub created_at : time::OffsetDateTime,pub updated_at : time::OffsetDateTime,}
impl<'a> From<GetSubmissionByIdBorrowed<'a>> for GetSubmissionById
{
    fn from(GetSubmissionByIdBorrowed { id,user_id,language,source_code,status,created_at,updated_at,}: GetSubmissionByIdBorrowed<'a>) -> Self
    { Self { id,user_id,language,source_code: source_code.into(),status,created_at,updated_at,} }
}pub struct GetSubmissionByIdQuery<'a, C: GenericClient, T, const N: usize>
{
    client: &'a  C, params:
    [&'a (dyn postgres_types::ToSql + Sync); N], stmt: &'a mut
    cornucopia_async::private::Stmt, extractor: fn(&tokio_postgres::Row) -> GetSubmissionByIdBorrowed,
    mapper: fn(GetSubmissionByIdBorrowed) -> T,
} impl<'a, C, T:'a, const N: usize> GetSubmissionByIdQuery<'a, C, T, N> where C:
GenericClient
{
    pub fn map<R>(self, mapper: fn(GetSubmissionByIdBorrowed) -> R) ->
    GetSubmissionByIdQuery<'a,C,R,N>
    {
        GetSubmissionByIdQuery
        {
            client: self.client, params: self.params, stmt: self.stmt,
            extractor: self.extractor, mapper,
        }
    } pub async fn one(self) -> Result<T, tokio_postgres::Error>
    {
        let stmt = self.stmt.prepare(self.client).await?; let row =
        self.client.query_one(stmt, &self.params).await?;
        Ok((self.mapper)((self.extractor)(&row)))
    } pub async fn all(self) -> Result<Vec<T>, tokio_postgres::Error>
    { self.iter().await?.try_collect().await } pub async fn opt(self) ->
    Result<Option<T>, tokio_postgres::Error>
    {
        let stmt = self.stmt.prepare(self.client).await?;
        Ok(self.client.query_opt(stmt, &self.params) .await?
        .map(|row| (self.mapper)((self.extractor)(&row))))
    } pub async fn iter(self,) -> Result<impl futures::Stream<Item = Result<T,
    tokio_postgres::Error>> + 'a, tokio_postgres::Error>
    {
        let stmt = self.stmt.prepare(self.client).await?; let it =
        self.client.query_raw(stmt,
        cornucopia_async::private::slice_iter(&self.params)) .await?
        .map(move |res|
        res.map(|row| (self.mapper)((self.extractor)(&row)))) .into_stream();
        Ok(it)
    }
}pub fn get_submission_by_id() -> GetSubmissionByIdStmt
{ GetSubmissionByIdStmt(cornucopia_async::private::Stmt::new("SELECT id, user_id, language, source_code, status, created_at, updated_at
FROM submissions
WHERE id = $1")) } pub struct
GetSubmissionByIdStmt(cornucopia_async::private::Stmt); impl GetSubmissionByIdStmt
{ pub fn bind<'a, C:
GenericClient,>(&'a mut self, client: &'a  C,
id: &'a i32,) -> GetSubmissionByIdQuery<'a,C, GetSubmissionById,
1>
{
    GetSubmissionByIdQuery
    {
        client, params: [id,], stmt: &mut self.0, extractor:
        |row| { GetSubmissionByIdBorrowed { id: row.get(0),user_id: row.get(1),language: row.get(2),source_code: row.get(3),status: row.get(4),created_at: row.get(5),updated_at: row.get(6),} }, mapper: |it| { <GetSubmissionById>::from(it) },
    }
} }}}