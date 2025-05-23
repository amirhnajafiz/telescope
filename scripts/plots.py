import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns


df = pd.read_csv('./analysis/datasets/final_dataset.csv')

filtered_df = df[df['Proxy Replication'] == 1].copy()

filtered_df['Caching'] = filtered_df['Caching'].replace({
    'None': 'None',
    'in-memory': 'In-memory',
    'file': 'File',
    'In-memory': 'In-memory',
    'File': 'File',
    None: 'None'
})

grouped = filtered_df.groupby(['Caching', 'Smart ABR'])['QoE'].mean().reset_index()

plt.figure(figsize=(10, 6))
sns.barplot(data=grouped, x='Caching', y='QoE', hue='Smart ABR',
            order=['None', 'In-memory', 'File'],
            palette=['#3476b7', '#c64141'])

plt.title('QoE by Caching Level and ABR Strategy')
plt.xlabel('Caching Level')
plt.ylabel('Average QoE')
plt.legend(title='ABR Strategy')
plt.tight_layout()

plt.savefig('qoe_by_caching_abr_simple.pyplot.png')
plt.show()


# Plot 1: QoE vs IPFS Bandwidth for each ABR Type 

plt.figure(figsize=(12, 6))
sns.scatterplot(data=df, x='IPFS Bandwidth (Mbps)', y='QoE', hue='Smart ABR')

plt.title('QoE vs IPFS Bandwidth for each ABR Type')
plt.xlabel('IPFS Bandwidth (Mbps)')
plt.ylabel('QoE')
plt.legend(title='Smart ABR')
plt.tight_layout()

plt.savefig('qoe_vs_ipfs_bandwidth_by_abr.pyplot.png')
plt.show()


# Plot 2: QoE by Caching Level and ABR Strategy (Proxy Replication = 1, Ordered)

filtered_df = df[df['Proxy Replication'] == 1].copy()

filtered_df['Caching'] = filtered_df['Caching'].replace({
    'None': 'None',
    'in-memory': 'In-memory',
    'file': 'File',
    'In-memory': 'In-memory',
    'File': 'File',
    None: 'None'
})

grouped_qoe = filtered_df.groupby(['Caching', 'Smart ABR'])['QoE'].mean().reset_index()

plt.figure(figsize=(10, 6))
sns.barplot(data=grouped_qoe, x='Caching', y='QoE', hue='Smart ABR',
            order=['None', 'In-memory', 'File'],
            palette=['#3476b7', '#c64141'])

plt.title('QoE by Caching Level and ABR Strategy (Proxy Replication = 1, Ordered)')
plt.xlabel('Caching Level')
plt.ylabel('Average QoE')
plt.legend(title='ABR Strategy')
plt.tight_layout()

plt.savefig('qoe_by_caching_abr_proxy_1_ordered.pyplot.png')
plt.show()

# Plot 3: QoE by Caching Level and ABR Strategy (Proxy Replication = 2, Ordered)

filtered_df_2 = df[df['Proxy Replication'] == 2].copy()

filtered_df_2['Caching'] = filtered_df_2['Caching'].replace({
    'None': 'None',
    'in-memory': 'In-memory',
    'file': 'File',
    'In-memory': 'In-memory',
    'File': 'File',
    None: 'None'
})

grouped_qoe_2 = filtered_df_2.groupby(['Caching', 'Smart ABR'])['QoE'].mean().reset_index()

plt.figure(figsize=(10, 6))
sns.barplot(data=grouped_qoe_2, x='Caching', y='QoE', hue='Smart ABR',
            order=['None', 'In-memory', 'File'],
            palette=['#3476b7', '#c64141'])

plt.title('QoE by Caching Level and ABR Strategy (Proxy Replication = 2, Ordered)')
plt.xlabel('Caching Level')
plt.ylabel('Average QoE')
plt.legend(title='ABR Strategy')
plt.tight_layout()

plt.savefig('qoe_by_caching_abr_proxy_2_ordered.pyplot.png')
plt.show()

# Plot 4: QoE by Caching Level and ABR Strategy (Proxy Replication = 3, Ordered)

filtered_df_3 = df[df['Proxy Replication'] == 3].copy()

filtered_df_3['Caching'] = filtered_df_3['Caching'].replace({
    'None': 'None',
    'in-memory': 'In-memory',
    'file': 'File',
    'In-memory': 'In-memory',
    'File': 'File',
    None: 'None'
})

grouped_qoe_3 = filtered_df_3.groupby(['Caching', 'Smart ABR'])['QoE'].mean().reset_index()

plt.figure(figsize=(10, 6))
sns.barplot(data=grouped_qoe_3, x='Caching', y='QoE', hue='Smart ABR',
            order=['None', 'In-memory', 'File'],
            palette=['#3476b7', '#c64141'])

plt.title('QoE by Caching Level and ABR Strategy (Proxy Replication = 3, Ordered)')
plt.xlabel('Caching Level')
plt.ylabel('Average QoE')
plt.legend(title='ABR Strategy')
plt.tight_layout()

plt.savefig('qoe_by_caching_abr_proxy_3_ordered.pyplot.png')
plt.show()

# Plot 5: Average Gateway Bandwidth per ABR Type

avg_gateway_bw = df.groupby('Smart ABR')['Gateway Bandwidth (Mbps)'].mean().reset_index()

plt.figure(figsize=(10, 6))
sns.barplot(data=avg_gateway_bw, x='Smart ABR', y='Gateway Bandwidth (Mbps)', color='orange')

plt.title('Average Gateway Bandwidth per ABR Type')
plt.xlabel('ABR Algorithm')
plt.ylabel('Gateway Bandwidth (Mbps)')
plt.tight_layout()

plt.savefig('average_gateway_bandwidth_by_abr.pyplot.png')
plt.show()

# Plot 6: Average IPFS Bandwidth per ABR Type

import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns

df = pd.read_csv('final_dataset.csv')

abr_bandwidth = df.groupby('Smart ABR')['IPFS Bandwidth (Mbps)'].mean().reset_index()

plt.figure(figsize=(10, 6))
sns.barplot(data=abr_bandwidth, x='Smart ABR', y='IPFS Bandwidth (Mbps)', color='orange')

plt.title('Average IPFS Bandwidth per ABR Type')
plt.xlabel('ABR Algorithm')
plt.ylabel('IPFS Bandwidth (Mbps)')
plt.tight_layout()

plt.savefig('average_ipfs_bandwidth_by_abr.pyplot.png')
plt.show()


# Plot 7: Average Client Stall Rate by Caching Level and ABR Strategy 

df['Caching'] = df['Caching'].replace({
    'None': 'None',
    'in-memory': 'In-memory',
    'file': 'File',
    'In-memory': 'In-memory',
    'File': 'File',
    None: 'None'
})

stall_rate_avg = df.groupby(['Caching', 'Smart ABR'])['Client Stall Rate'].mean().reset_index()

plt.figure(figsize=(10, 6))
sns.barplot(data=stall_rate_avg, x='Caching', y='Client Stall Rate', hue='Smart ABR',
            order=['None', 'In-memory', 'File'],
            palette=['#3476b7', '#c64141'])

plt.title('Average Client Stall Rate by Caching Level and ABR Strategy')
plt.xlabel('Caching Level')
plt.ylabel('Average Stall Rate')
plt.legend(title='ABR Strategy')
plt.tight_layout()

plt.savefig('avg_stall_rate_by_caching_and_abr.pyplot.png')
plt.show()

# Plot 8: Average Client Stall Rate by Caching Level and ABR Strategy (Proxy Replication = 1) 

filtered_df_stall = df[df['Proxy Replication'] == 1].copy()

filtered_df_stall['Caching'] = filtered_df_stall['Caching'].replace({
    'None': 'None',
    'in-memory': 'In-memory',
    'file': 'File',
    'In-memory': 'In-memory',
    'File': 'File',
    None: 'None'
})

stall_data = filtered_df_stall.groupby(['Caching', 'Smart ABR'])['Client Stall Rate'].mean().reset_index()

plt.figure(figsize=(10, 6))
sns.barplot(data=stall_data, x='Caching', y='Client Stall Rate', hue='Smart ABR',
            order=['None', 'In-memory', 'File'],
            palette=['#3476b7', '#c64141'])

plt.title('Average Client Stall Rate by Caching Level and ABR Strategy (Proxy Replication = 1)')
plt.xlabel('Caching Level')
plt.ylabel('Average Stall Rate')
plt.legend(title='ABR Strategy')
plt.tight_layout()

plt.savefig('avg_stall_rate_by_caching_abr_proxy_1.pyplot.png')
plt.show()

# Plot 9: Average Client Stall Rate by Caching Level and ABR Strategy (Proxy Replication = 2)

filtered_df_stall_2 = df[df['Proxy Replication'] == 2].copy()

filtered_df_stall_2['Caching'] = filtered_df_stall_2['Caching'].replace({
    'None': 'None',
    'in-memory': 'In-memory',
    'file': 'File',
    'In-memory': 'In-memory',
    'File': 'File',
    None: 'None'
})

stall_data_2 = filtered_df_stall_2.groupby(['Caching', 'Smart ABR'])['Client Stall Rate'].mean().reset_index()

plt.figure(figsize=(10, 6))
sns.barplot(data=stall_data_2, x='Caching', y='Client Stall Rate', hue='Smart ABR',
            order=['None', 'In-memory', 'File'],
            palette=['#3476b7', '#c64141'])

plt.title('Average Client Stall Rate by Caching Level and ABR Strategy (Proxy Replication = 2)')
plt.xlabel('Caching Level')
plt.ylabel('Average Stall Rate')
plt.legend(title='ABR Strategy')
plt.tight_layout()

plt.savefig('avg_stall_rate_by_caching_abr_proxy_2.pyplot.png')
plt.show()

# Plot 10: Average Client Stall Rate by Caching Level and ABR Strategy (Proxy Replication = 3)

filtered_df_stall_3 = df[df['Proxy Replication'] == 3].copy()

filtered_df_stall_3['Caching'] = filtered_df_stall_3['Caching'].replace({
    'None': 'None',
    'in-memory': 'In-memory',
    'file': 'File',
    'In-memory': 'In-memory',
    'File': 'File',
    None: 'None'
})

stall_data_3 = filtered_df_stall_3.groupby(['Caching', 'Smart ABR'])['Client Stall Rate'].mean().reset_index()

plt.figure(figsize=(10, 6))
sns.barplot(data=stall_data_3, x='Caching', y='Client Stall Rate', hue='Smart ABR',
            order=['None', 'In-memory', 'File'],
            palette=['#3476b7', '#c64141'])

plt.title('Average Client Stall Rate by Caching Level and ABR Strategy (Proxy Replication = 3)')
plt.xlabel('Caching Level')
plt.ylabel('Average Stall Rate')
plt.legend(title='ABR Strategy')
plt.tight_layout()

plt.savefig('avg_stall_rate_by_caching_abr_proxy_3.pyplot.png')
plt.show()

# Plot 11: QoE vs Video Resolution (Ordered, by ABR Strategy)

quality_mapping = {
    1: '640x360',
    2: '854x480',
    3: '1280x720',
    4: '4K'
}

df['Video Resolution'] = df['Client Video Quality'].map(quality_mapping)
df['Video Resolution'] = pd.Categorical(df['Video Resolution'], categories=['640x360', '854x480', '1280x720', '4K'], ordered=True)

plt.figure(figsize=(10, 6))
sns.lineplot(data=df, x='Video Resolution', y='QoE', hue='Smart ABR',
             hue_order=['statistics-based', 'throughput-based'],
             markers=True, style='Smart ABR', ci='sd')

plt.title('QoE vs Video Resolution (640x360 to 4K, by ABR Strategy)')
plt.xlabel('Video Resolution')
plt.ylabel('QoE')
plt.tight_layout()

plt.savefig('qoe_vs_video_resolution_by_abr_ordered.pyplot.png')
plt.show()


# Plot 12: QoE vs IPFS Hops (by Caching Type)

df['Caching'] = df['Caching'].replace({
    'none': 'None',
    'in-memory': 'In-memory',
    'file': 'File',
    'In-memory': 'In-memory',
    'File': 'File',
    'None': 'None'
})

plt.figure(figsize=(10, 6))
sns.lineplot(data=df, x='IPFS Hops', y='QoE', hue='Caching',
             hue_order=['None', 'In-memory', 'File'],
             markers=True, style='Caching', ci='sd')

plt.title('QoE vs IPFS Hops (by Caching Type)')
plt.xlabel('IPFS Hops')
plt.ylabel('QoE')
plt.tight_layout()

plt.savefig('qoe_vs_ipfs_hops_by_caching.pyplot.png')
plt.show()


# Plot 13: Correlation Heatmap (Filtered)

df_corr = df.drop(columns=['Test Index', 'Number of Tests'], errors='ignore')
numerical_df = df_corr.select_dtypes(include=['number'])
corr_matrix = numerical_df.corr()

plt.figure(figsize=(12, 10))
sns.heatmap(corr_matrix, annot=True, cmap='coolwarm', center=0, fmt=".2f")
plt.title('Correlation Heatmap of Numerical Metrics')
plt.tight_layout()
plt.savefig('correlation_heatmap_filtered.pyplot.png')
plt.show()


ordered_video_names = [
    'big_buck_bunny_720p_5mb',
    'big_buck_bunny_720p_10mb',
    'big_buck_bunny_720p_20mb',
    'big_buck_bunny_720p_30mb'
]


heatmap_data = df.pivot_table(index='Video Name', columns='Caching', values='QoE', aggfunc='mean')

heatmap_data = heatmap_data.reindex(ordered_video_names)
heatmap_data = heatmap_data[['File', 'In-memory', 'None']]

# Heatmap
plt.figure(figsize=(10, 8))
sns.heatmap(heatmap_data, annot=True, cmap='RdBu_r', fmt=".2f")
plt.title('QoE by Video Name and Caching Strategy')
plt.xlabel('Caching Strategy')
plt.ylabel('Video Name')
plt.tight_layout()
plt.savefig('qoe_by_video_and_caching_heatmap_ordered.pyplot.png')
plt.show()

# Plot 14: QoE vs. Average Stall Rate by Caching Strategy

df['Video Size'] = df['Video Name'].str.extract(r'_(\d+mb)')
df['Caching'] = df['Caching'].replace({
    'file': 'File',
    'in-memory': 'In-memory',
    'none': 'None',
    'File': 'File',
    'In-memory': 'In-memory',
    'None': 'None'
})

grouped_scatter = df.groupby(['Caching', 'Video Size'])[['QoE', 'Client Stall Rate']].mean().reset_index()
grouped_scatter.rename(columns={'Client Stall Rate': 'Average Stall Rate'}, inplace=True)

plt.figure(figsize=(12, 7))
sns.scatterplot(data=grouped_scatter,
                x='Average Stall Rate', y='QoE',
                hue='Caching', style='Video Size', s=100)

plt.title('QoE vs. Average Stall Rate by Caching Strategy')
plt.xlabel('Average Stall Rate')
plt.ylabel('QoE')
plt.tight_layout()
plt.savefig('qoe_vs_stall_rate_by_caching_strategy.pyplot.png')
plt.show()
