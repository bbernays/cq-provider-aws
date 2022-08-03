package kinesis

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/firehose"
	"github.com/aws/aws-sdk-go-v2/service/firehose/types"
	"github.com/cloudquery/cq-provider-aws/client"
	"github.com/cloudquery/cq-provider-sdk/provider/diag"
	"github.com/cloudquery/cq-provider-sdk/provider/schema"
)

func Firehoses() *schema.Table {
	return &schema.Table{
		Name:         "aws_kinesis_firehoses",
		Description:  "Contains information about a delivery stream",
		Resolver:     fetchKinesisFirehoses,
		Multiplex:    client.ServiceAccountRegionMultiplexer("firehose"),
		IgnoreError:  client.IgnoreAccessDeniedServiceDisabled,
		DeleteFilter: client.DeleteAccountRegionFilter,
		Options:      schema.TableCreationOptions{PrimaryKeys: []string{"arn"}},
		Columns: []schema.Column{
			{
				Name:        "account_id",
				Description: "The AWS Account ID of the resource.",
				Type:        schema.TypeString,
				Resolver:    client.ResolveAWSAccount,
			},
			{
				Name:        "region",
				Description: "The AWS Region of the resource.",
				Type:        schema.TypeString,
				Resolver:    client.ResolveAWSRegion,
			},
			{
				Name:     "tags",
				Type:     schema.TypeJSON,
				Resolver: resolveKinesisFirehoseTags,
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the delivery stream",
				Type:        schema.TypeString,
				Resolver:    schema.PathResolver("DeliveryStreamARN"),
			},
			{
				Name:        "delivery_stream_arn",
				Description: "The Amazon Resource Name (ARN) of the delivery stream",
				Type:        schema.TypeString,
				Resolver:    schema.PathResolver("DeliveryStreamARN"),
			},
			{
				Name:        "delivery_stream_name",
				Description: "The name of the delivery stream",
				Type:        schema.TypeString,
			},
			{
				Name:        "delivery_stream_status",
				Description: "The status of the delivery stream",
				Type:        schema.TypeString,
			},
			{
				Name:        "delivery_stream_type",
				Description: "The delivery stream type",
				Type:        schema.TypeString,
			},
			{
				Name:        "version_id",
				Description: "Each time the destination is updated for a delivery stream, the version ID is changed, and the current version ID is required when updating the destination This is so that the service knows it is applying the changes to the correct version of the delivery stream",
				Type:        schema.TypeString,
			},
			{
				Name:        "create_timestamp",
				Description: "The date and time that the delivery stream was created",
				Type:        schema.TypeTimestamp,
			},
			{
				Name:        "encryption_config_failure_description_details",
				Description: "A message providing details about the error that caused the failure",
				Type:        schema.TypeString,
				Resolver:    schema.PathResolver("DeliveryStreamEncryptionConfiguration.FailureDescription.Details"),
			},
			{
				Name:        "encryption_config_failure_description_type",
				Description: "The type of error that caused the failure",
				Type:        schema.TypeString,
				Resolver:    schema.PathResolver("DeliveryStreamEncryptionConfiguration.FailureDescription.Type"),
			},
			{
				Name:        "encryption_config_key_arn",
				Description: "If KeyType is CUSTOMER_MANAGED_CMK, this field contains the ARN of the customer managed CMK",
				Type:        schema.TypeString,
				Resolver:    schema.PathResolver("DeliveryStreamEncryptionConfiguration.KeyARN"),
			},
			{
				Name:        "encryption_config_key_type",
				Description: "Indicates the type of customer master key (CMK) that is used for encryption",
				Type:        schema.TypeString,
				Resolver:    schema.PathResolver("DeliveryStreamEncryptionConfiguration.KeyType"),
			},
			{
				Name:        "encryption_config_status",
				Description: "This is the server-side encryption (SSE) status for the delivery stream",
				Type:        schema.TypeString,
				Resolver:    schema.PathResolver("DeliveryStreamEncryptionConfiguration.Status"),
			},
			{
				Name:        "failure_description_details",
				Description: "A message providing details about the error that caused the failure",
				Type:        schema.TypeString,
				Resolver:    schema.PathResolver("FailureDescription.Details"),
			},
			{
				Name:        "failure_description_type",
				Description: "The type of error that caused the failure",
				Type:        schema.TypeString,
				Resolver:    schema.PathResolver("FailureDescription.Type"),
			},
			{
				Name:        "last_update_timestamp",
				Description: "The date and time that the delivery stream was last updated",
				Type:        schema.TypeTimestamp,
			},
			{
				Name:        "source_kinesis_stream_delivery_start_timestamp",
				Description: "Kinesis Data Firehose starts retrieving records from the Kinesis data stream starting with this timestamp",
				Type:        schema.TypeTimestamp,
				Resolver:    schema.PathResolver("Source.KinesisStreamSourceDescription.DeliveryStartTimestamp"),
			},
			{
				Name:        "source_kinesis_stream_kinesis_stream_arn",
				Description: "The Amazon Resource Name (ARN) of the source Kinesis data stream",
				Type:        schema.TypeString,
				Resolver:    schema.PathResolver("Source.KinesisStreamSourceDescription.KinesisStreamARN"),
			},
			{
				Name:        "source_kinesis_stream_role_arn",
				Description: "The ARN of the role used by the source Kinesis data stream",
				Type:        schema.TypeString,
				Resolver:    schema.PathResolver("Source.KinesisStreamSourceDescription.RoleARN"),
			},
		},
		Relations: []*schema.Table{
			{
				Name:     "aws_kinesis_firehose_open_search_destination",
				Resolver: schema.PathTableResolver("Destinations.AmazonopensearchserviceDestinationDescription"),
				Columns: []schema.Column{
					{
						Name:        "firehose_cq_id",
						Description: "Unique CloudQuery ID of aws_kinesis_firehoses table (FK)",
						Type:        schema.TypeUUID,
						Resolver:    schema.ParentIdResolver,
					},
					{
						Name:        "processing_configuration_processors",
						Description: "The Amazon Resource Name (ARN) of the delivery stream",
						Type:        schema.TypeJSON,
						Resolver:    resolveKinesisFirehoseOpenSearchDestinationProcessingConfigurationProcessors,
					},
					{
						Name:     "buffering_hints_interval_in_seconds",
						Type:     schema.TypeBigInt,
						Resolver: schema.PathResolver("BufferingHints.IntervalInSeconds"),
					},
					{
						Name:     "buffering_hints_size_in_mb_s",
						Type:     schema.TypeBigInt,
						Resolver: schema.PathResolver("BufferingHints.SizeInMBs"),
					},
					{
						Name:        "cloud_watch_logging_options_enabled",
						Description: "Enables or disables CloudWatch logging",
						Type:        schema.TypeBool,
						Resolver:    schema.PathResolver("CloudWatchLoggingOptions.Enabled"),
					},
					{
						Name:        "cloud_watch_logging_options_log_group_name",
						Description: "The CloudWatch group name for logging",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("CloudWatchLoggingOptions.LogGroupName"),
					},
					{
						Name:        "cloud_watch_logging_options_log_stream_name",
						Description: "The CloudWatch log stream name for logging",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("CloudWatchLoggingOptions.LogStreamName"),
					},
					{
						Name: "cluster_endpoint",
						Type: schema.TypeString,
					},
					{
						Name:     "domain_arn",
						Type:     schema.TypeString,
						Resolver: schema.PathResolver("DomainARN"),
					},
					{
						Name: "index_name",
						Type: schema.TypeString,
					},
					{
						Name: "index_rotation_period",
						Type: schema.TypeString,
					},
					{
						Name:        "processing_configuration_enabled",
						Description: "Enables or disables data processing",
						Type:        schema.TypeBool,
						Resolver:    schema.PathResolver("ProcessingConfiguration.Enabled"),
					},
					{
						Name:     "retry_options_duration_in_seconds",
						Type:     schema.TypeBigInt,
						Resolver: schema.PathResolver("RetryOptions.DurationInSeconds"),
					},
					{
						Name:     "role_arn",
						Type:     schema.TypeString,
						Resolver: schema.PathResolver("RoleARN"),
					},
					{
						Name: "s3_backup_mode",
						Type: schema.TypeString,
					},
					{
						Name:        "s3_destination_bucket_arn",
						Description: "The ARN of the S3 bucket",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("S3DestinationDescription.BucketARN"),
					},
					{
						Name:        "s3_destination_buffering_hints_interval_in_seconds",
						Description: "Buffer incoming data for the specified period of time, in seconds, before delivering it to the destination",
						Type:        schema.TypeBigInt,
						Resolver:    schema.PathResolver("S3DestinationDescription.BufferingHints.IntervalInSeconds"),
					},
					{
						Name:        "s3_destination_buffering_hints_size_in_mb_s",
						Description: "Buffer incoming data to the specified size, in MiBs, before delivering it to the destination",
						Type:        schema.TypeBigInt,
						Resolver:    schema.PathResolver("S3DestinationDescription.BufferingHints.SizeInMBs"),
					},
					{
						Name:        "s3_destination_compression_format",
						Description: "The compression format",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("S3DestinationDescription.CompressionFormat"),
					},
					{
						Name:        "s3_destination_kms_encryption_config_aws_kms_key_arn",
						Description: "The Amazon Resource Name (ARN) of the encryption key",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("S3DestinationDescription.EncryptionConfiguration.KMSEncryptionConfig.AWSKMSKeyARN"),
					},
					{
						Name:        "s3_destination_no_encryption_config",
						Description: "Specifically override existing encryption information to ensure that no encryption is used",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("S3DestinationDescription.EncryptionConfiguration.NoEncryptionConfig"),
					},
					{
						Name:        "s3_destination_role_arn",
						Description: "The Amazon Resource Name (ARN) of the AWS credentials",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("S3DestinationDescription.RoleARN"),
					},
					{
						Name:        "s3_destination_cloud_watch_logging_options_enabled",
						Description: "Enables or disables CloudWatch logging",
						Type:        schema.TypeBool,
						Resolver:    schema.PathResolver("S3DestinationDescription.CloudWatchLoggingOptions.Enabled"),
					},
					{
						Name:        "s3_destination_cloud_watch_logging_options_log_group_name",
						Description: "The CloudWatch group name for logging",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("S3DestinationDescription.CloudWatchLoggingOptions.LogGroupName"),
					},
					{
						Name:        "s3_destination_cloud_watch_logging_options_log_stream_name",
						Description: "The CloudWatch log stream name for logging",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("S3DestinationDescription.CloudWatchLoggingOptions.LogStreamName"),
					},
					{
						Name:        "s3_destination_error_output_prefix",
						Description: "A prefix that Kinesis Data Firehose evaluates and adds to failed records before writing them to S3",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("S3DestinationDescription.ErrorOutputPrefix"),
					},
					{
						Name:        "s3_destination_prefix",
						Description: "The \"YYYY/MM/DD/HH\" time format prefix is automatically used for delivered Amazon S3 files",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("S3DestinationDescription.Prefix"),
					},
					{
						Name: "type_name",
						Type: schema.TypeString,
					},
					{
						Name:        "vpc_configuration_description_role_arn",
						Description: "The ARN of the IAM role that the delivery stream uses to create endpoints in the destination VPC",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("VpcConfigurationDescription.RoleARN"),
					},
					{
						Name:        "vpc_configuration_description_security_group_ids",
						Description: "The IDs of the security groups that Kinesis Data Firehose uses when it creates ENIs in the VPC of the Amazon ES destination",
						Type:        schema.TypeStringArray,
						Resolver:    schema.PathResolver("VpcConfigurationDescription.SecurityGroupIds"),
					},
					{
						Name:        "vpc_configuration_description_subnet_ids",
						Description: "The IDs of the subnets that Kinesis Data Firehose uses to create ENIs in the VPC of the Amazon ES destination",
						Type:        schema.TypeStringArray,
						Resolver:    schema.PathResolver("VpcConfigurationDescription.SubnetIds"),
					},
					{
						Name:        "vpc_configuration_description_vpc_id",
						Description: "The ID of the Amazon ES destination's VPC",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("VpcConfigurationDescription.VpcId"),
					},
				},
			},
			{
				Name:        "aws_kinesis_firehose_extended_s3_destination",
				Description: "Describes a destination in Amazon S3",
				Resolver:    schema.PathTableResolver("Destinations.ExtendedS3DestinationDescription"),
				Columns: []schema.Column{
					{
						Name:        "firehose_cq_id",
						Description: "Unique CloudQuery ID of aws_kinesis_firehoses table (FK)",
						Type:        schema.TypeUUID,
						Resolver:    schema.ParentIdResolver,
					},
					{
						Name:        "processing_configuration_processors",
						Description: "The Amazon Resource Name (ARN) of the delivery stream",
						Type:        schema.TypeJSON,
						Resolver:    resolveKinesisFirehoseExtendedS3DestinationProcessingConfigurationProcessors,
					},
					{
						Name:        "bucket_arn",
						Description: "The ARN of the S3 bucket",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("BucketARN"),
					},
					{
						Name:        "buffering_hints_interval_in_seconds",
						Description: "Buffer incoming data for the specified period of time, in seconds, before delivering it to the destination",
						Type:        schema.TypeBigInt,
						Resolver:    schema.PathResolver("BufferingHints.IntervalInSeconds"),
					},
					{
						Name:        "buffering_hints_size_in_mb_s",
						Description: "Buffer incoming data to the specified size, in MiBs, before delivering it to the destination",
						Type:        schema.TypeBigInt,
						Resolver:    schema.PathResolver("BufferingHints.SizeInMBs"),
					},
					{
						Name:        "compression_format",
						Description: "The compression format",
						Type:        schema.TypeString,
					},
					{
						Name:        "encryption_configuration_kms_encryption_config_aws_kms_key_arn",
						Description: "The Amazon Resource Name (ARN) of the encryption key",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("EncryptionConfiguration.KMSEncryptionConfig.AWSKMSKeyARN"),
					},
					{
						Name:        "encryption_configuration_no_encryption_config",
						Description: "Specifically override existing encryption information to ensure that no encryption is used",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("EncryptionConfiguration.NoEncryptionConfig"),
					},
					{
						Name:        "role_arn",
						Description: "The Amazon Resource Name (ARN) of the AWS credentials",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("RoleARN"),
					},
					{
						Name:        "cloud_watch_logging_options_enabled",
						Description: "Enables or disables CloudWatch logging",
						Type:        schema.TypeBool,
						Resolver:    schema.PathResolver("CloudWatchLoggingOptions.Enabled"),
					},
					{
						Name:        "cloud_watch_logging_options_log_group_name",
						Description: "The CloudWatch group name for logging",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("CloudWatchLoggingOptions.LogGroupName"),
					},
					{
						Name:        "cloud_watch_logging_options_log_stream_name",
						Description: "The CloudWatch log stream name for logging",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("CloudWatchLoggingOptions.LogStreamName"),
					},
					{
						Name:        "enabled",
						Description: "Defaults to true",
						Type:        schema.TypeBool,
						Resolver:    schema.PathResolver("DataFormatConversionConfiguration.Enabled"),
					},
					{
						Name:        "deserializer_hive_json_ser_de_timestamp_formats",
						Description: "Indicates how you want Kinesis Data Firehose to parse the date and timestamps that may be present in your input data JSON",
						Type:        schema.TypeStringArray,
						Resolver:    schema.PathResolver("DataFormatConversionConfiguration.InputFormatConfiguration.Deserializer.HiveJsonSerDe.TimestampFormats"),
					},
					{
						Name:        "deserializer_open_x_json_ser_de_case_insensitive",
						Description: "When set to true, which is the default, Kinesis Data Firehose converts JSON keys to lowercase before deserializing them",
						Type:        schema.TypeBool,
						Resolver:    schema.PathResolver("DataFormatConversionConfiguration.InputFormatConfiguration.Deserializer.OpenXJsonSerDe.CaseInsensitive"),
					},
					{
						Name:        "deserializer_open_x_json_ser_de_column_to_json_key_mappings",
						Description: "Maps column names to JSON keys that aren't identical to the column names",
						Type:        schema.TypeJSON,
						Resolver:    schema.PathResolver("DataFormatConversionConfiguration.InputFormatConfiguration.Deserializer.OpenXJsonSerDe.ColumnToJsonKeyMappings"),
					},
					{
						Name:        "serializer_orc_ser_de_block_size_bytes",
						Description: "The Hadoop Distributed File System (HDFS) block size",
						Type:        schema.TypeBigInt,
						Resolver:    schema.PathResolver("DataFormatConversionConfiguration.OutputFormatConfiguration.Serializer.OrcSerDe.BlockSizeBytes"),
					},
					{
						Name:        "serializer_orc_ser_de_bloom_filter_columns",
						Description: "The column names for which you want Kinesis Data Firehose to create bloom filters",
						Type:        schema.TypeStringArray,
						Resolver:    schema.PathResolver("DataFormatConversionConfiguration.OutputFormatConfiguration.Serializer.OrcSerDe.BloomFilterColumns"),
					},
					{
						Name:        "serializer_orc_ser_de_bloom_filter_false_positive_probability",
						Description: "The Bloom filter false positive probability (FPP)",
						Type:        schema.TypeFloat,
						Resolver:    schema.PathResolver("DataFormatConversionConfiguration.OutputFormatConfiguration.Serializer.OrcSerDe.BloomFilterFalsePositiveProbability"),
					},
					{
						Name:        "serializer_orc_ser_de_compression",
						Description: "The compression code to use over data blocks",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("DataFormatConversionConfiguration.OutputFormatConfiguration.Serializer.OrcSerDe.Compression"),
					},
					{
						Name:        "serializer_orc_ser_de_dictionary_key_threshold",
						Description: "Represents the fraction of the total number of non-null rows",
						Type:        schema.TypeFloat,
						Resolver:    schema.PathResolver("DataFormatConversionConfiguration.OutputFormatConfiguration.Serializer.OrcSerDe.DictionaryKeyThreshold"),
					},
					{
						Name:        "serializer_orc_ser_de_enable_padding",
						Description: "Set this to true to indicate that you want stripes to be padded to the HDFS block boundaries",
						Type:        schema.TypeBool,
						Resolver:    schema.PathResolver("DataFormatConversionConfiguration.OutputFormatConfiguration.Serializer.OrcSerDe.EnablePadding"),
					},
					{
						Name:        "serializer_orc_ser_de_format_version",
						Description: "The version of the file to write",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("DataFormatConversionConfiguration.OutputFormatConfiguration.Serializer.OrcSerDe.FormatVersion"),
					},
					{
						Name:        "serializer_orc_ser_de_padding_tolerance",
						Description: "A number between 0 and 1 that defines the tolerance for block padding as a decimal fraction of stripe size",
						Type:        schema.TypeFloat,
						Resolver:    schema.PathResolver("DataFormatConversionConfiguration.OutputFormatConfiguration.Serializer.OrcSerDe.PaddingTolerance"),
					},
					{
						Name:        "serializer_orc_ser_de_row_index_stride",
						Description: "The number of rows between index entries",
						Type:        schema.TypeBigInt,
						Resolver:    schema.PathResolver("DataFormatConversionConfiguration.OutputFormatConfiguration.Serializer.OrcSerDe.RowIndexStride"),
					},
					{
						Name:        "serializer_orc_ser_de_stripe_size_bytes",
						Description: "The number of bytes in each stripe",
						Type:        schema.TypeBigInt,
						Resolver:    schema.PathResolver("DataFormatConversionConfiguration.OutputFormatConfiguration.Serializer.OrcSerDe.StripeSizeBytes"),
					},
					{
						Name:        "serializer_parquet_ser_de_block_size_bytes",
						Description: "The Hadoop Distributed File System (HDFS) block size",
						Type:        schema.TypeBigInt,
						Resolver:    schema.PathResolver("DataFormatConversionConfiguration.OutputFormatConfiguration.Serializer.ParquetSerDe.BlockSizeBytes"),
					},
					{
						Name:        "serializer_parquet_ser_de_compression",
						Description: "The compression code to use over data blocks",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("DataFormatConversionConfiguration.OutputFormatConfiguration.Serializer.ParquetSerDe.Compression"),
					},
					{
						Name:        "serializer_parquet_ser_de_enable_dictionary_compression",
						Description: "Indicates whether to enable dictionary compression",
						Type:        schema.TypeBool,
						Resolver:    schema.PathResolver("DataFormatConversionConfiguration.OutputFormatConfiguration.Serializer.ParquetSerDe.EnableDictionaryCompression"),
					},
					{
						Name:        "serializer_parquet_ser_de_max_padding_bytes",
						Description: "The maximum amount of padding to apply",
						Type:        schema.TypeBigInt,
						Resolver:    schema.PathResolver("DataFormatConversionConfiguration.OutputFormatConfiguration.Serializer.ParquetSerDe.MaxPaddingBytes"),
					},
					{
						Name:        "serializer_parquet_ser_de_page_size_bytes",
						Description: "The Parquet page size",
						Type:        schema.TypeBigInt,
						Resolver:    schema.PathResolver("DataFormatConversionConfiguration.OutputFormatConfiguration.Serializer.ParquetSerDe.PageSizeBytes"),
					},
					{
						Name:        "serializer_parquet_ser_de_writer_version",
						Description: "Indicates the version of row format to output",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("DataFormatConversionConfiguration.OutputFormatConfiguration.Serializer.ParquetSerDe.WriterVersion"),
					},
					{
						Name:        "schema_configuration_catalog_id",
						Description: "The ID of the AWS Glue Data Catalog",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("DataFormatConversionConfiguration.SchemaConfiguration.CatalogId"),
					},
					{
						Name:        "schema_configuration_database_name",
						Description: "Specifies the name of the AWS Glue database that contains the schema for the output data",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("DataFormatConversionConfiguration.SchemaConfiguration.DatabaseName"),
					},
					{
						Name:        "schema_configuration_region",
						Description: "If you don't specify an AWS Region, the default is the current Region",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("DataFormatConversionConfiguration.SchemaConfiguration.Region"),
					},
					{
						Name:        "schema_configuration_role_arn",
						Description: "The role that Kinesis Data Firehose can use to access AWS Glue",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("DataFormatConversionConfiguration.SchemaConfiguration.RoleARN"),
					},
					{
						Name:        "schema_configuration_table_name",
						Description: "Specifies the AWS Glue table that contains the column information that constitutes your data schema",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("DataFormatConversionConfiguration.SchemaConfiguration.TableName"),
					},
					{
						Name:        "schema_configuration_version_id",
						Description: "Specifies the table version for the output data schema",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("DataFormatConversionConfiguration.SchemaConfiguration.VersionId"),
					},
					{
						Name:        "dynamic_partitioning_configuration_enabled",
						Description: "Specifies that the dynamic partitioning is enabled for this Kinesis Data Firehose delivery stream",
						Type:        schema.TypeBool,
						Resolver:    schema.PathResolver("DynamicPartitioningConfiguration.Enabled"),
					},
					{
						Name:        "error_output_prefix",
						Description: "A prefix that Kinesis Data Firehose evaluates and adds to failed records before writing them to S3",
						Type:        schema.TypeString,
					},
					{
						Name:        "prefix",
						Description: "The \"YYYY/MM/DD/HH\" time format prefix is automatically used for delivered Amazon S3 files",
						Type:        schema.TypeString,
					},
					{
						Name:        "processing_configuration_enabled",
						Description: "Enables or disables data processing",
						Type:        schema.TypeBool,
						Resolver:    schema.PathResolver("ProcessingConfiguration.Enabled"),
					},
					{
						Name:        "s3_backup_bucket_arn",
						Description: "The ARN of the S3 bucket",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("S3BackupDescription.BucketARN"),
					},
					{
						Name:        "s3_backup_buffering_hints_interval_in_seconds",
						Description: "Buffer incoming data for the specified period of time, in seconds, before delivering it to the destination",
						Type:        schema.TypeBigInt,
						Resolver:    schema.PathResolver("S3BackupDescription.BufferingHints.IntervalInSeconds"),
					},
					{
						Name:        "s3_backup_buffering_hints_size_in_mb_s",
						Description: "Buffer incoming data to the specified size, in MiBs, before delivering it to the destination",
						Type:        schema.TypeBigInt,
						Resolver:    schema.PathResolver("S3BackupDescription.BufferingHints.SizeInMBs"),
					},
					{
						Name:        "s3_backup_compression_format",
						Description: "The compression format",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("S3BackupDescription.CompressionFormat"),
					},
					{
						Name:        "s3_backup_kms_encryption_config_aws_kms_key_arn",
						Description: "The Amazon Resource Name (ARN) of the encryption key",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("S3BackupDescription.EncryptionConfiguration.KMSEncryptionConfig.AWSKMSKeyARN"),
					},
					{
						Name:        "s3_backup_no_encryption_config",
						Description: "Specifically override existing encryption information to ensure that no encryption is used",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("S3BackupDescription.EncryptionConfiguration.NoEncryptionConfig"),
					},
					{
						Name:        "s3_backup_role_arn",
						Description: "The Amazon Resource Name (ARN) of the AWS credentials",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("S3BackupDescription.RoleARN"),
					},
					{
						Name:        "s3_backup_cloud_watch_logging_options_enabled",
						Description: "Enables or disables CloudWatch logging",
						Type:        schema.TypeBool,
						Resolver:    schema.PathResolver("S3BackupDescription.CloudWatchLoggingOptions.Enabled"),
					},
					{
						Name:        "s3_backup_cloud_watch_logging_options_log_group_name",
						Description: "The CloudWatch group name for logging",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("S3BackupDescription.CloudWatchLoggingOptions.LogGroupName"),
					},
					{
						Name:        "s3_backup_cloud_watch_logging_options_log_stream_name",
						Description: "The CloudWatch log stream name for logging",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("S3BackupDescription.CloudWatchLoggingOptions.LogStreamName"),
					},
					{
						Name:        "s3_backup_error_output_prefix",
						Description: "A prefix that Kinesis Data Firehose evaluates and adds to failed records before writing them to S3",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("S3BackupDescription.ErrorOutputPrefix"),
					},
					{
						Name:        "s3_backup_prefix",
						Description: "The \"YYYY/MM/DD/HH\" time format prefix is automatically used for delivered Amazon S3 files",
						Type:        schema.TypeString,
						Resolver:    schema.PathResolver("S3BackupDescription.Prefix"),
					},
					{
						Name:        "s3_backup_mode",
						Description: "The Amazon S3 backup mode",
						Type:        schema.TypeString,
					},
				},
			},
		},
	}
}

// ====================================================================================================================
//                                               Table Resolver Functions
// ====================================================================================================================

func fetchKinesisFirehoses(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	return diag.WrapError(client.ListAndDetailResolver(ctx, meta, res, listDeliveryStreams, deliveryStreamDetail))
}
func resolveKinesisFirehoseTags(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
	cl := meta.(*client.Client)
	svc := cl.Services().Firehose
	summary := resource.Item.(*types.DeliveryStreamDescription)
	input := firehose.ListTagsForDeliveryStreamInput{
		DeliveryStreamName: summary.DeliveryStreamName,
	}
	var tags []types.Tag
	for {
		output, err := svc.ListTagsForDeliveryStream(ctx, &input)
		if err != nil {
			return diag.WrapError(err)
		}
		tags = append(tags, output.Tags...)
		if !aws.ToBool(output.HasMoreTags) {
			break
		}
		input.ExclusiveStartTagKey = aws.String(*output.Tags[len(output.Tags)-1].Key)
	}
	return diag.WrapError(resource.Set(c.Name, client.TagsToMap(tags)))
}
func resolveKinesisFirehoseOpenSearchDestinationProcessingConfigurationProcessors(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
	r := resource.Item.(*types.AmazonopensearchserviceDestinationDescription)
	out, err := json.Marshal(r.ProcessingConfiguration.Processors)
	if err != nil {
		return diag.WrapError(err)
	}
	return diag.WrapError(resource.Set(c.Name, out))
}
func resolveKinesisFirehoseExtendedS3DestinationProcessingConfigurationProcessors(ctx context.Context, meta schema.ClientMeta, resource *schema.Resource, c schema.Column) error {
	r := resource.Item.(*types.ExtendedS3DestinationDescription)
	out, err := json.Marshal(r.ProcessingConfiguration.Processors)
	if err != nil {
		return diag.WrapError(err)
	}
	return diag.WrapError(resource.Set(c.Name, out))
}

// ====================================================================================================================
//                                                  User Defined Helpers
// ====================================================================================================================

func listDeliveryStreams(ctx context.Context, meta schema.ClientMeta, detailChan chan<- interface{}) error {
	c := meta.(*client.Client)
	svc := c.Services().Firehose
	input := firehose.ListDeliveryStreamsInput{}
	for {
		response, err := svc.ListDeliveryStreams(ctx, &input)
		if err != nil {
			return diag.WrapError(err)
		}
		for _, item := range response.DeliveryStreamNames {
			detailChan <- item
		}
		if !aws.ToBool(response.HasMoreDeliveryStreams) {
			break
		}
		input.ExclusiveStartDeliveryStreamName = aws.String(response.DeliveryStreamNames[len(response.DeliveryStreamNames)-1])
	}
	return nil
}
func deliveryStreamDetail(ctx context.Context, meta schema.ClientMeta, resultsChan chan<- interface{}, errorChan chan<- error, listInfo interface{}) {
	c := meta.(*client.Client)
	streamName := listInfo.(string)
	svc := c.Services().Firehose
	streamSummary, err := svc.DescribeDeliveryStream(ctx, &firehose.DescribeDeliveryStreamInput{
		DeliveryStreamName: aws.String(streamName),
	})
	if err != nil {
		if c.IsNotFoundError(err) {
			return
		}
		errorChan <- diag.WrapError(err)
		return
	}
	resultsChan <- streamSummary.DeliveryStreamDescription
}
